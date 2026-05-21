// Package guard implements the stateless circuit breaker middleware chain
// that intercepts every tool call before execution. Each Guard in the chain
// is a pure function of the request — no disk I/O, no DB, no session affinity.
package guard

import (
	"fmt"
	"time"
)

// ToolRequest represents a single tool invocation request to be validated.
type ToolRequest struct {
	ToolName   string         `json:"tool_name"`
	Args       map[string]any `json:"args"`
	SessionID  string         `json:"session_id"`
	TurnNumber int            `json:"turn_number"`
	Timestamp  time.Time      `json:"timestamp"`
}

// ToolResponse is the verdict from a Guard check.
type ToolResponse struct {
	Allowed   bool   `json:"allowed"`
	Reason    string `json:"reason,omitempty"`
	GuardName string `json:"guard_name,omitempty"` // which guard blocked it
}

// Guard is a single middleware in the chain. Implementations must be
// safe for concurrent use — use atomic operations or sync primitives
// for any mutable state (counters, circuit state).
type Guard interface {
	// Check evaluates whether a tool request should be allowed.
	// Must be fast (<1ms) and allocation-free on the hot path.
	Check(req ToolRequest) ToolResponse

	// Name returns a human-readable identifier for this guard.
	Name() string
}

// PostExecutionHook is called after tool execution to feed results
// back into guards that need them (e.g., circuit breaker tracking failures).
type PostExecutionHook interface {
	RecordResult(toolName string, err error)
}

// Breaker chains multiple guards into a single middleware pipeline.
// It is safe for concurrent use.
type Breaker struct {
	guards  []Guard
	metrics *Metrics
}

// NewBreaker creates a new Breaker with the given guards evaluated in order.
func NewBreaker(guards ...Guard) *Breaker {
	return &Breaker{
		guards:  guards,
		metrics: NewMetrics(),
	}
}

// Execute runs the full guard chain against a tool request.
// Returns on the first rejection. If all guards pass, returns Allowed: true.
func (b *Breaker) Execute(req ToolRequest) ToolResponse {
	for _, g := range b.guards {
		resp := g.Check(req)
		if !resp.Allowed {
			resp.GuardName = g.Name()
			b.metrics.RecordRejection(g.Name())
			return resp
		}
	}
	b.metrics.RecordAllow()
	return ToolResponse{Allowed: true}
}

// RecordResult propagates tool execution results to any guards that
// implement PostExecutionHook (e.g., CircuitBreakerGuard).
func (b *Breaker) RecordResult(toolName string, err error) {
	for _, g := range b.guards {
		if hook, ok := g.(PostExecutionHook); ok {
			hook.RecordResult(toolName, err)
		}
	}
}

// Stats returns a snapshot of guard metrics.
func (b *Breaker) Stats() MetricsSnapshot {
	return b.metrics.Snapshot()
}

// String returns a human-readable summary of the guard chain.
func (b *Breaker) String() string {
	names := make([]string, len(b.guards))
	for i, g := range b.guards {
		names[i] = g.Name()
	}
	return fmt.Sprintf("Breaker[%v]", names)
}
