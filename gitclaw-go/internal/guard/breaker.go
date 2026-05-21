package guard

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// CircuitState represents the state of a circuit breaker for a specific tool.
type CircuitState int32

const (
	// StateClosed is the normal operating state — all calls allowed.
	StateClosed CircuitState = iota

	// StateOpen means the circuit is tripped — all calls rejected.
	StateOpen

	// StateHalfOpen allows a single probe call to test recovery.
	StateHalfOpen
)

func (s CircuitState) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// toolCircuit holds the per-tool circuit breaker state.
type toolCircuit struct {
	state    atomic.Int32 // CircuitState
	failures atomic.Int32
	openedAt atomic.Int64 // unix nano
}

func (tc *toolCircuit) getState() CircuitState {
	return CircuitState(tc.state.Load())
}

func (tc *toolCircuit) setState(s CircuitState) {
	tc.state.Store(int32(s))
}

func (tc *toolCircuit) getOpenedAt() time.Time {
	return time.Unix(0, tc.openedAt.Load())
}

// CircuitBreakerGuard trips open after consecutive failures and rejects
// calls until a reset timeout elapses. After timeout, allows one probe
// call (half-open). If the probe succeeds, closes the circuit.
type CircuitBreakerGuard struct {
	failureThreshold int32
	resetTimeout     time.Duration
	circuits         sync.Map // toolName → *toolCircuit
}

// CircuitBreakerConfig configures the circuit breaker.
type CircuitBreakerConfig struct {
	// FailureThreshold is the number of consecutive failures to trip the circuit.
	FailureThreshold int `yaml:"failure_threshold" json:"failure_threshold"`

	// ResetTimeout is how long the circuit stays open before allowing a probe.
	ResetTimeout time.Duration `yaml:"reset_timeout" json:"reset_timeout"`
}

// NewCircuitBreakerGuard creates a new CircuitBreakerGuard.
func NewCircuitBreakerGuard(cfg CircuitBreakerConfig) *CircuitBreakerGuard {
	if cfg.FailureThreshold == 0 {
		cfg.FailureThreshold = 5
	}
	if cfg.ResetTimeout == 0 {
		cfg.ResetTimeout = 30 * time.Second
	}
	return &CircuitBreakerGuard{
		failureThreshold: int32(cfg.FailureThreshold),
		resetTimeout:     cfg.ResetTimeout,
	}
}

func (c *CircuitBreakerGuard) Name() string { return "circuit_breaker" }

func (c *CircuitBreakerGuard) Check(req ToolRequest) ToolResponse {
	tc := c.getOrCreateCircuit(req.ToolName)

	switch tc.getState() {
	case StateOpen:
		// Check if reset timeout has elapsed
		if time.Since(tc.getOpenedAt()) > c.resetTimeout {
			tc.setState(StateHalfOpen)
			// Allow one probe call
			return ToolResponse{Allowed: true}
		}
		return ToolResponse{
			Allowed: false,
			Reason: fmt.Sprintf(
				"circuit breaker OPEN for tool %q (%d consecutive failures, resets in %v)",
				req.ToolName,
				tc.failures.Load(),
				c.resetTimeout-time.Since(tc.getOpenedAt()),
			),
		}

	case StateHalfOpen:
		// Already probing — reject additional calls while probe is in flight
		return ToolResponse{
			Allowed: false,
			Reason:  fmt.Sprintf("circuit breaker HALF-OPEN for tool %q (probe in progress)", req.ToolName),
		}

	default: // StateClosed
		return ToolResponse{Allowed: true}
	}
}

// RecordResult implements PostExecutionHook. Called after each tool execution
// to update failure counts and circuit state.
func (c *CircuitBreakerGuard) RecordResult(toolName string, err error) {
	tc := c.getOrCreateCircuit(toolName)

	if err != nil {
		count := tc.failures.Add(1)
		if count >= c.failureThreshold {
			tc.setState(StateOpen)
			tc.openedAt.Store(time.Now().UnixNano())
		}
	} else {
		// Success — reset failures and close circuit
		tc.failures.Store(0)
		tc.setState(StateClosed)
	}
}

// GetState returns the current circuit state for a tool (for diagnostics).
func (c *CircuitBreakerGuard) GetState(toolName string) CircuitState {
	tc := c.getOrCreateCircuit(toolName)
	return tc.getState()
}

func (c *CircuitBreakerGuard) getOrCreateCircuit(toolName string) *toolCircuit {
	if v, ok := c.circuits.Load(toolName); ok {
		return v.(*toolCircuit)
	}
	tc := &toolCircuit{}
	actual, _ := c.circuits.LoadOrStore(toolName, tc)
	return actual.(*toolCircuit)
}
