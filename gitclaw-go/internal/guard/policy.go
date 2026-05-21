package guard

import (
	"fmt"
	"regexp"
	"strings"
)

// PolicyRule defines a single allow/deny rule for tool invocations.
type PolicyRule struct {
	// Tool is the tool name to match. Use "*" for all tools.
	Tool string `yaml:"tool" json:"tool"`

	// Action is "allow" or "deny".
	Action string `yaml:"action,omitempty" json:"action,omitempty"`

	// Deny is a shorthand pattern — if set, Action is implicitly "deny".
	// Pattern format: "args.<key> matches '<regex>'" or "args contains '<substring>'"
	Deny string `yaml:"deny,omitempty" json:"deny,omitempty"`

	// Allow is a shorthand pattern for explicit allow rules.
	Allow string `yaml:"allow,omitempty" json:"allow,omitempty"`

	// Description explains what this rule does (for audit logs).
	Description string `yaml:"description,omitempty" json:"description,omitempty"`

	// compiled regex cache (populated on first use)
	compiledRegex *regexp.Regexp
	patternParsed bool
	patternField  string // e.g., "command"
	patternType   string // "matches" or "contains"
	patternValue  string // the regex or substring
}

// parsePattern extracts field, match type, and value from pattern strings like:
// "args.command matches 'rm -rf.*'"
// "args contains 'sudo'"
func (r *PolicyRule) parsePattern() {
	if r.patternParsed {
		return
	}
	r.patternParsed = true

	pattern := r.Deny
	if pattern == "" {
		pattern = r.Allow
	}
	if pattern == "" {
		return
	}

	// Pattern: args.<field> matches '<regex>'
	if strings.HasPrefix(pattern, "args.") {
		parts := strings.SplitN(pattern, " ", 3)
		if len(parts) >= 3 {
			r.patternField = strings.TrimPrefix(parts[0], "args.")
			r.patternType = parts[1]
			r.patternValue = strings.Trim(parts[2], "'\"")
			if r.patternType == "matches" {
				r.compiledRegex, _ = regexp.Compile(r.patternValue)
			}
		}
		return
	}

	// Pattern: args contains '<substring>'
	if strings.HasPrefix(pattern, "args contains ") {
		r.patternType = "contains"
		r.patternValue = strings.Trim(strings.TrimPrefix(pattern, "args contains "), "'\"")
		return
	}
}

// matches checks if a tool request matches this rule.
func (r *PolicyRule) matches(req ToolRequest) bool {
	// Check tool name match
	if r.Tool != "*" && r.Tool != req.ToolName {
		return false
	}

	r.parsePattern()

	// No pattern — matches all calls to this tool
	if r.patternType == "" {
		return true
	}

	switch r.patternType {
	case "matches":
		if r.patternField != "" {
			// Check specific arg field
			val, ok := req.Args[r.patternField]
			if !ok {
				return false
			}
			strVal := fmt.Sprintf("%v", val)
			if r.compiledRegex != nil {
				return r.compiledRegex.MatchString(strVal)
			}
			return false
		}

	case "contains":
		if r.patternField != "" {
			// Check specific arg field
			val, ok := req.Args[r.patternField]
			if !ok {
				return false
			}
			return strings.Contains(fmt.Sprintf("%v", val), r.patternValue)
		}
		// Check all arg values
		for _, v := range req.Args {
			if strings.Contains(fmt.Sprintf("%v", v), r.patternValue) {
				return true
			}
		}
		return false
	}

	return false
}

// PolicyChecker evaluates tool requests against a list of policy rules.
// Rules are evaluated in order; first matching deny rule wins.
type PolicyChecker struct {
	rules []PolicyRule
}

// PolicyCheckerConfig holds the policy rules.
type PolicyCheckerConfig struct {
	Rules []PolicyRule `yaml:"policies" json:"policies"`
}

// NewPolicyChecker creates a new PolicyChecker from the given config.
func NewPolicyChecker(cfg PolicyCheckerConfig) *PolicyChecker {
	return &PolicyChecker{rules: cfg.Rules}
}

func (p *PolicyChecker) Name() string { return "policy_checker" }

func (p *PolicyChecker) Check(req ToolRequest) ToolResponse {
	for i := range p.rules {
		rule := &p.rules[i]

		isDeny := rule.Deny != "" || rule.Action == "deny"
		if !isDeny {
			continue
		}

		if rule.matches(req) {
			reason := rule.Description
			if reason == "" {
				pattern := rule.Deny
				if pattern == "" {
					pattern = "tool blocked by policy"
				}
				reason = fmt.Sprintf("policy violation: %s", pattern)
			}
			return ToolResponse{
				Allowed: false,
				Reason:  reason,
			}
		}
	}

	return ToolResponse{Allowed: true}
}
