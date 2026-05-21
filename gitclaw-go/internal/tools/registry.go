// Package tools provides the built-in tool implementations and a registry
// for managing all available tools (built-in + declarative + plugin).
package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/open-gitagent/gitclaw-go/internal/agent"
	"github.com/open-gitagent/gitclaw-go/internal/state"
)

// ToolConfig holds configuration shared by all built-in tools.
type ToolConfig struct {
	AgentDir   string
	Timeout    int // seconds
	Ledger     *state.Ledger
	SessionID  string
}

// CreateBuiltinTools returns all built-in tools wired with the given config.
func CreateBuiltinTools(cfg ToolConfig) []*agent.ToolDef {
	return []*agent.ToolDef{
		NewCLITool(cfg),
		NewReadTool(cfg),
		NewWriteTool(cfg),
		NewMemoryTool(cfg),
	}
}

// ToolDefsToLLM converts agent ToolDefs to LLM API tool definitions.
func ToolDefsToLLM(tools []*agent.ToolDef) []agent.LLMToolDef {
	defs := make([]agent.LLMToolDef, len(tools))
	for i, t := range tools {
		defs[i] = agent.LLMToolDef{
			Type: "function",
			Function: agent.LLMFunctionDef{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.InputSchema,
			},
		}
	}
	return defs
}

// getStringArg safely extracts a string argument from the args map.
func getStringArg(args map[string]any, key string) string {
	if v, ok := args[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
		return fmt.Sprintf("%v", v)
	}
	return ""
}

// getBoolArg safely extracts a boolean argument from the args map.
func getBoolArg(args map[string]any, key string) bool {
	if v, ok := args[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// withTimeout creates a context with the configured timeout.
func withTimeout(ctx context.Context, timeoutSeconds int) (context.Context, context.CancelFunc) {
	if timeoutSeconds <= 0 {
		timeoutSeconds = 120
	}
	return context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
}
