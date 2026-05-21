package guard

import (
	"fmt"
	"math"
	"sync/atomic"
)

// CostGuard enforces hard spending and token ceilings. When the cumulative
// cost or token count exceeds the configured maximum, all further tool
// calls are rejected until the session ends.
type CostGuard struct {
	maxCostUsd float64
	maxTokens  int64

	// currentCostMicro stores cost in micro-dollars (1 USD = 1_000_000)
	// to avoid floating point atomics.
	currentCostMicro atomic.Int64
	currentTokens    atomic.Int64
}

// CostGuardConfig configures the cost guard.
type CostGuardConfig struct {
	// MaxCostUsd is the maximum allowed cost in USD (e.g., 10.00).
	MaxCostUsd float64 `yaml:"max_usd" json:"max_usd"`

	// MaxTokens is the maximum total tokens (input + output) allowed.
	MaxTokens int64 `yaml:"max_tokens" json:"max_tokens"`
}

// NewCostGuard creates a new CostGuard.
func NewCostGuard(cfg CostGuardConfig) *CostGuard {
	if cfg.MaxCostUsd == 0 {
		cfg.MaxCostUsd = math.MaxFloat64
	}
	if cfg.MaxTokens == 0 {
		cfg.MaxTokens = math.MaxInt64
	}
	return &CostGuard{
		maxCostUsd: cfg.MaxCostUsd,
		maxTokens:  cfg.MaxTokens,
	}
}

func (c *CostGuard) Name() string { return "cost_guard" }

func (c *CostGuard) Check(req ToolRequest) ToolResponse {
	costUsd := c.CurrentCostUsd()
	if costUsd >= c.maxCostUsd {
		return ToolResponse{
			Allowed: false,
			Reason: fmt.Sprintf(
				"cost ceiling exceeded: $%.4f / $%.2f",
				costUsd, c.maxCostUsd,
			),
		}
	}

	tokens := c.currentTokens.Load()
	if tokens >= c.maxTokens {
		return ToolResponse{
			Allowed: false,
			Reason: fmt.Sprintf(
				"token ceiling exceeded: %d / %d",
				tokens, c.maxTokens,
			),
		}
	}

	return ToolResponse{Allowed: true}
}

// AddCost records a cost increment. costUsd is in dollars (e.g., 0.003).
func (c *CostGuard) AddCost(costUsd float64) {
	microDollars := int64(costUsd * 1_000_000)
	c.currentCostMicro.Add(microDollars)
}

// AddTokens records a token count increment.
func (c *CostGuard) AddTokens(count int64) {
	c.currentTokens.Add(count)
}

// CurrentCostUsd returns the current accumulated cost in USD.
func (c *CostGuard) CurrentCostUsd() float64 {
	return float64(c.currentCostMicro.Load()) / 1_000_000.0
}

// CurrentTokens returns the current accumulated token count.
func (c *CostGuard) CurrentTokens() int64 {
	return c.currentTokens.Load()
}

// Remaining returns the remaining budget and tokens before the guard trips.
func (c *CostGuard) Remaining() (costUsd float64, tokens int64) {
	costUsd = c.maxCostUsd - c.CurrentCostUsd()
	if costUsd < 0 {
		costUsd = 0
	}
	tokens = c.maxTokens - c.currentTokens.Load()
	if tokens < 0 {
		tokens = 0
	}
	return
}
