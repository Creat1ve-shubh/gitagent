package guard_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/open-gitagent/gitclaw-go/internal/guard"
)

// ── Middleware / Breaker tests ─────────────────────────────────────────

func TestBreakerAllowsWhenNoGuards(t *testing.T) {
	b := guard.NewBreaker()
	resp := b.Execute(guard.ToolRequest{ToolName: "cli", SessionID: "s1"})
	if !resp.Allowed {
		t.Fatalf("expected allowed, got blocked: %s", resp.Reason)
	}
}

func TestBreakerRejectsOnFirstFailingGuard(t *testing.T) {
	deny := &alwaysDeny{name: "test_deny"}
	allow := &alwaysAllow{name: "test_allow"}

	// deny first → should block
	b := guard.NewBreaker(deny, allow)
	resp := b.Execute(guard.ToolRequest{ToolName: "cli"})
	if resp.Allowed {
		t.Fatal("expected rejection")
	}
	if resp.GuardName != "test_deny" {
		t.Fatalf("expected guard name 'test_deny', got %q", resp.GuardName)
	}
}

func TestBreakerMetrics(t *testing.T) {
	deny := &alwaysDeny{name: "test_deny"}
	b := guard.NewBreaker(deny)

	for i := 0; i < 5; i++ {
		b.Execute(guard.ToolRequest{ToolName: "cli"})
	}

	stats := b.Stats()
	if stats.TotalRejected != 5 {
		t.Fatalf("expected 5 rejections, got %d", stats.TotalRejected)
	}
	if stats.TotalAllowed != 0 {
		t.Fatalf("expected 0 allowed, got %d", stats.TotalAllowed)
	}
}

// ── Rate limiter tests ─────────────────────────────────────────────────

func TestRateLimiterAllowsWithinLimit(t *testing.T) {
	rl := guard.NewRateLimiter(guard.RateLimiterConfig{
		WindowSize:   10 * time.Second,
		MaxPerWindow: 5,
	})

	for i := 0; i < 5; i++ {
		resp := rl.Check(guard.ToolRequest{SessionID: "s1", ToolName: "cli"})
		if !resp.Allowed {
			t.Fatalf("call %d should be allowed, got: %s", i+1, resp.Reason)
		}
	}
}

func TestRateLimiterBlocksOverLimit(t *testing.T) {
	rl := guard.NewRateLimiter(guard.RateLimiterConfig{
		WindowSize:   10 * time.Second,
		MaxPerWindow: 3,
	})

	for i := 0; i < 3; i++ {
		rl.Check(guard.ToolRequest{SessionID: "s1", ToolName: "cli"})
	}

	resp := rl.Check(guard.ToolRequest{SessionID: "s1", ToolName: "cli"})
	if resp.Allowed {
		t.Fatal("4th call should be blocked")
	}
}

func TestRateLimiterIsolatesSessions(t *testing.T) {
	rl := guard.NewRateLimiter(guard.RateLimiterConfig{
		WindowSize:   10 * time.Second,
		MaxPerWindow: 2,
	})

	// Exhaust session 1
	rl.Check(guard.ToolRequest{SessionID: "s1", ToolName: "cli"})
	rl.Check(guard.ToolRequest{SessionID: "s1", ToolName: "cli"})

	// Session 2 should still be allowed
	resp := rl.Check(guard.ToolRequest{SessionID: "s2", ToolName: "cli"})
	if !resp.Allowed {
		t.Fatal("session 2 should not be rate limited by session 1")
	}
}

func TestRateLimiterConcurrentSafety(t *testing.T) {
	rl := guard.NewRateLimiter(guard.RateLimiterConfig{
		WindowSize:   10 * time.Second,
		MaxPerWindow: 1000,
	})

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				rl.Check(guard.ToolRequest{
					SessionID: fmt.Sprintf("s%d", n%5),
					ToolName:  "cli",
				})
			}
		}(i)
	}
	wg.Wait()
	// If we get here without a race condition, the test passes
}

// ── Policy checker tests ───────────────────────────────────────────────

func TestPolicyBlocksRmRf(t *testing.T) {
	pc := guard.NewPolicyChecker(guard.PolicyCheckerConfig{
		Rules: []guard.PolicyRule{
			{Tool: "cli", Deny: "args.command matches 'rm -rf.*'"},
		},
	})

	resp := pc.Check(guard.ToolRequest{
		ToolName: "cli",
		Args:     map[string]any{"command": "rm -rf /"},
	})
	if resp.Allowed {
		t.Fatal("rm -rf should be blocked")
	}

	resp = pc.Check(guard.ToolRequest{
		ToolName: "cli",
		Args:     map[string]any{"command": "ls -la"},
	})
	if !resp.Allowed {
		t.Fatalf("ls should be allowed, got: %s", resp.Reason)
	}
}

func TestPolicyBlocksSudo(t *testing.T) {
	pc := guard.NewPolicyChecker(guard.PolicyCheckerConfig{
		Rules: []guard.PolicyRule{
			{Tool: "*", Deny: "args contains 'sudo'"},
		},
	})

	resp := pc.Check(guard.ToolRequest{
		ToolName: "cli",
		Args:     map[string]any{"command": "sudo apt install nginx"},
	})
	if resp.Allowed {
		t.Fatal("sudo should be blocked")
	}
}

func TestPolicyAllowsSafeCli(t *testing.T) {
	pc := guard.NewPolicyChecker(guard.PolicyCheckerConfig{
		Rules: []guard.PolicyRule{
			{Tool: "cli", Deny: "args.command matches 'rm -rf.*'"},
		},
	})

	resp := pc.Check(guard.ToolRequest{
		ToolName: "cli",
		Args:     map[string]any{"command": "git status"},
	})
	if !resp.Allowed {
		t.Fatalf("git status should be allowed, got: %s", resp.Reason)
	}
}

func TestPolicyIgnoresNonMatchingTools(t *testing.T) {
	pc := guard.NewPolicyChecker(guard.PolicyCheckerConfig{
		Rules: []guard.PolicyRule{
			{Tool: "cli", Deny: "args.command matches 'rm -rf.*'"},
		},
	})

	resp := pc.Check(guard.ToolRequest{
		ToolName: "write", // not "cli"
		Args:     map[string]any{"command": "rm -rf /"},
	})
	if !resp.Allowed {
		t.Fatal("write tool should not be affected by cli policy")
	}
}

// ── Circuit breaker tests ──────────────────────────────────────────────

func TestCircuitBreakerStartsClosed(t *testing.T) {
	cb := guard.NewCircuitBreakerGuard(guard.CircuitBreakerConfig{
		FailureThreshold: 3,
		ResetTimeout:     1 * time.Second,
	})

	resp := cb.Check(guard.ToolRequest{ToolName: "cli"})
	if !resp.Allowed {
		t.Fatal("circuit should start closed (allowing calls)")
	}
}

func TestCircuitBreakerTripsAfterThreshold(t *testing.T) {
	cb := guard.NewCircuitBreakerGuard(guard.CircuitBreakerConfig{
		FailureThreshold: 3,
		ResetTimeout:     5 * time.Second,
	})

	// Record 3 failures
	for i := 0; i < 3; i++ {
		cb.RecordResult("cli", fmt.Errorf("tool failed"))
	}

	// Should be open now
	resp := cb.Check(guard.ToolRequest{ToolName: "cli"})
	if resp.Allowed {
		t.Fatal("circuit should be open after 3 failures")
	}

	st := cb.GetState("cli")
	if st.String() != "open" {
		t.Fatalf("expected state 'open', got %q", st)
	}
}

func TestCircuitBreakerResetsOnSuccess(t *testing.T) {
	cb := guard.NewCircuitBreakerGuard(guard.CircuitBreakerConfig{
		FailureThreshold: 3,
		ResetTimeout:     5 * time.Second,
	})

	// 2 failures, then a success
	cb.RecordResult("cli", fmt.Errorf("fail"))
	cb.RecordResult("cli", fmt.Errorf("fail"))
	cb.RecordResult("cli", nil) // success resets

	// Should still be closed
	resp := cb.Check(guard.ToolRequest{ToolName: "cli"})
	if !resp.Allowed {
		t.Fatal("circuit should be closed after success")
	}
}

func TestCircuitBreakerIsolatesTools(t *testing.T) {
	cb := guard.NewCircuitBreakerGuard(guard.CircuitBreakerConfig{
		FailureThreshold: 2,
		ResetTimeout:     5 * time.Second,
	})

	// Trip circuit for "cli"
	cb.RecordResult("cli", fmt.Errorf("fail"))
	cb.RecordResult("cli", fmt.Errorf("fail"))

	// "read" should still be allowed
	resp := cb.Check(guard.ToolRequest{ToolName: "read"})
	if !resp.Allowed {
		t.Fatal("read tool circuit should be independent of cli")
	}
}

// ── Cost guard tests ───────────────────────────────────────────────────

func TestCostGuardAllowsWithinBudget(t *testing.T) {
	cg := guard.NewCostGuard(guard.CostGuardConfig{
		MaxCostUsd: 1.0,
		MaxTokens:  10000,
	})

	cg.AddCost(0.5)
	cg.AddTokens(5000)

	resp := cg.Check(guard.ToolRequest{ToolName: "cli"})
	if !resp.Allowed {
		t.Fatalf("should be within budget, got: %s", resp.Reason)
	}
}

func TestCostGuardBlocksOverBudget(t *testing.T) {
	cg := guard.NewCostGuard(guard.CostGuardConfig{
		MaxCostUsd: 1.0,
		MaxTokens:  10000,
	})

	cg.AddCost(1.01)

	resp := cg.Check(guard.ToolRequest{ToolName: "cli"})
	if resp.Allowed {
		t.Fatal("should block when over cost budget")
	}
}

func TestCostGuardBlocksOverTokens(t *testing.T) {
	cg := guard.NewCostGuard(guard.CostGuardConfig{
		MaxCostUsd: 100.0,
		MaxTokens:  1000,
	})

	cg.AddTokens(1001)

	resp := cg.Check(guard.ToolRequest{ToolName: "cli"})
	if resp.Allowed {
		t.Fatal("should block when over token limit")
	}
}

func TestCostGuardRemaining(t *testing.T) {
	cg := guard.NewCostGuard(guard.CostGuardConfig{
		MaxCostUsd: 10.0,
		MaxTokens:  50000,
	})

	cg.AddCost(3.5)
	cg.AddTokens(20000)

	costRemaining, tokensRemaining := cg.Remaining()
	if costRemaining < 6.4 || costRemaining > 6.6 {
		t.Fatalf("expected ~6.5 remaining, got %f", costRemaining)
	}
	if tokensRemaining != 30000 {
		t.Fatalf("expected 30000 tokens remaining, got %d", tokensRemaining)
	}
}

// ── Full pipeline integration test ─────────────────────────────────────

func TestFullPipelineIntegration(t *testing.T) {
	breaker := guard.BuildBreaker(guard.Config{
		RateLimit: &guard.RateLimiterConfig{
			WindowSize:   10 * time.Second,
			MaxPerWindow: 5,
		},
		CircuitBreaker: &guard.CircuitBreakerConfig{
			FailureThreshold: 3,
			ResetTimeout:     1 * time.Second,
		},
		CostCeiling: &guard.CostGuardConfig{
			MaxCostUsd: 10.0,
			MaxTokens:  100000,
		},
		Policies: []guard.PolicyRule{
			{Tool: "cli", Deny: "args.command matches 'rm -rf.*'"},
		},
	})

	// Normal call — should pass all guards
	resp := breaker.Execute(guard.ToolRequest{
		ToolName:  "cli",
		Args:      map[string]any{"command": "ls"},
		SessionID: "test",
		Timestamp: time.Now(),
	})
	if !resp.Allowed {
		t.Fatalf("normal call should be allowed: %s", resp.Reason)
	}

	// Dangerous call — should be blocked by policy
	resp = breaker.Execute(guard.ToolRequest{
		ToolName:  "cli",
		Args:      map[string]any{"command": "rm -rf /tmp/important"},
		SessionID: "test",
		Timestamp: time.Now(),
	})
	if resp.Allowed {
		t.Fatal("rm -rf should be blocked by policy")
	}
	if resp.GuardName != "policy_checker" {
		t.Fatalf("expected policy_checker to block, got %q", resp.GuardName)
	}

	stats := breaker.Stats()
	if stats.TotalAllowed != 1 || stats.TotalRejected != 1 {
		t.Fatalf("expected 1 allowed, 1 rejected; got %+v", stats)
	}
}

// ── Test helpers ───────────────────────────────────────────────────────

type alwaysDeny struct{ name string }

func (a *alwaysDeny) Check(_ guard.ToolRequest) guard.ToolResponse {
	return guard.ToolResponse{Allowed: false, Reason: "always deny"}
}
func (a *alwaysDeny) Name() string { return a.name }

type alwaysAllow struct{ name string }

func (a *alwaysAllow) Check(_ guard.ToolRequest) guard.ToolResponse {
	return guard.ToolResponse{Allowed: true}
}
func (a *alwaysAllow) Name() string { return a.name }
