package guard

import "time"

// Config holds the full guard pipeline configuration,
// typically parsed from the `runtime.guard` section of agent.yaml.
type Config struct {
	RateLimit      *RateLimiterConfig    `yaml:"rate_limit,omitempty" json:"rate_limit,omitempty"`
	CircuitBreaker *CircuitBreakerConfig `yaml:"circuit_breaker,omitempty" json:"circuit_breaker,omitempty"`
	CostCeiling    *CostGuardConfig      `yaml:"cost_ceiling,omitempty" json:"cost_ceiling,omitempty"`
	Policies       []PolicyRule          `yaml:"policies,omitempty" json:"policies,omitempty"`
}

// DefaultConfig returns a sensible default guard configuration.
func DefaultConfig() Config {
	return Config{
		RateLimit: &RateLimiterConfig{
			WindowSize:   60 * time.Second,
			MaxPerWindow: 100,
		},
		CircuitBreaker: &CircuitBreakerConfig{
			FailureThreshold: 5,
			ResetTimeout:     30 * time.Second,
		},
		CostCeiling: &CostGuardConfig{
			MaxCostUsd: 10.0,
			MaxTokens:  500_000,
		},
		Policies: []PolicyRule{
			{
				Tool:        "cli",
				Deny:        "args.command matches 'rm -rf.*'",
				Description: "Block recursive force-delete commands",
			},
			{
				Tool:        "*",
				Deny:        "args contains 'sudo'",
				Description: "Block sudo commands in all tools",
			},
		},
	}
}

// BuildBreaker constructs a Breaker from a Config.
// Guards are ordered: RateLimiter → PolicyChecker → CircuitBreaker → CostGuard.
// This ordering ensures cheap checks run first (rate limit is O(1) atomic)
// and expensive checks (policy regex) run only if rate limit passes.
func BuildBreaker(cfg Config) *Breaker {
	var guards []Guard

	// 1. Rate limiter — cheapest check, runs first
	if cfg.RateLimit != nil {
		guards = append(guards, NewRateLimiter(*cfg.RateLimit))
	}

	// 2. Policy checker — regex matching on args
	if len(cfg.Policies) > 0 {
		guards = append(guards, NewPolicyChecker(PolicyCheckerConfig{Rules: cfg.Policies}))
	}

	// 3. Circuit breaker — per-tool failure tracking
	if cfg.CircuitBreaker != nil {
		guards = append(guards, NewCircuitBreakerGuard(*cfg.CircuitBreaker))
	}

	// 4. Cost guard — budget ceiling
	if cfg.CostCeiling != nil {
		guards = append(guards, NewCostGuard(*cfg.CostCeiling))
	}

	return NewBreaker(guards...)
}
