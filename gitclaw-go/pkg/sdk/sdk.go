// Package sdk provides the public Go API for using gitclaw programmatically.
//
//	result := sdk.Query(sdk.Options{
//	    Prompt: "Build a REST API",
//	    Dir:    "/path/to/agent",
//	})
//	for msg := range result.Messages() {
//	    fmt.Println(msg.Content)
//	}
package sdk

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	agentpkg "github.com/open-gitagent/gitclaw-go/internal/agent"
	"github.com/open-gitagent/gitclaw-go/internal/config"
	"github.com/open-gitagent/gitclaw-go/internal/guard"
	"github.com/open-gitagent/gitclaw-go/internal/hooks"
	"github.com/open-gitagent/gitclaw-go/internal/state"
	"github.com/open-gitagent/gitclaw-go/internal/tools"
)

// Options configures a query execution.
type Options struct {
	// Prompt is the user message to send to the agent.
	Prompt string

	// Dir is the agent directory (must contain agent.yaml).
	Dir string

	// Model overrides the model from agent.yaml (e.g., "openai:gpt-4o").
	Model string

	// SystemPrompt replaces the default system prompt entirely.
	SystemPrompt string

	// SystemPromptSuffix appends to the default system prompt.
	SystemPromptSuffix string

	// MaxTurns limits the number of LLM ↔ tool turns.
	MaxTurns int

	// AllowedTools whitelists tool names. Empty = all tools.
	AllowedTools []string

	// DisallowedTools blacklists tool names.
	DisallowedTools []string

	// GuardConfig overrides the default guard pipeline config.
	GuardConfig *guard.Config

	// Context for cancellation. If nil, context.Background() is used.
	Context context.Context
}

// Message is a single event from the agent execution.
type Message struct {
	Type    string // "assistant", "tool_use", "tool_result", "system", "error", "guard_block"
	Content string
	Meta    map[string]any
}

// Result holds the output of a query execution.
type Result struct {
	messages []Message
	err      error
}

// Messages returns all collected messages.
func (r *Result) Messages() []Message { return r.messages }

// Error returns any error from the execution.
func (r *Result) Error() error { return r.err }

// Query runs the agent with the given options and collects all messages.
func Query(opts Options) *Result {
	result := &Result{}

	ctx := opts.Context
	if ctx == nil {
		ctx = context.Background()
	}

	agentDir := opts.Dir
	if agentDir == "" {
		agentDir, _ = os.Getwd()
	}
	agentDir, _ = filepath.Abs(agentDir)

	// Load manifest
	manifest, err := config.LoadManifest(agentDir)
	if err != nil {
		result.err = fmt.Errorf("loading agent: %w", err)
		return result
	}

	if opts.Model != "" {
		manifest.Model.Preferred = opts.Model
	}
	if opts.MaxTurns > 0 {
		manifest.Runtime.MaxTurns = opts.MaxTurns
	}

	// Guard
	guardCfg := guard.DefaultConfig()
	if opts.GuardConfig != nil {
		guardCfg = *opts.GuardConfig
	} else if manifest.Runtime.Guard != nil {
		guardCfg = *manifest.Runtime.Guard
	}
	breaker := guard.BuildBreaker(guardCfg)

	// Ledger
	ledger := state.NewLedger(agentDir)
	loopCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go state.CommitLoop(loopCtx, ledger)

	// System prompt
	systemPrompt := config.LoadIdentityFiles(agentDir)
	if opts.SystemPrompt != "" {
		systemPrompt = opts.SystemPrompt
	}
	if opts.SystemPromptSuffix != "" {
		systemPrompt += "\n\n" + opts.SystemPromptSuffix
	}

	// Tools
	sessionID := fmt.Sprintf("sdk-%d", os.Getpid())
	builtinTools := tools.CreateBuiltinTools(tools.ToolConfig{
		AgentDir:  agentDir,
		Timeout:   manifest.Runtime.Timeout,
		Ledger:    ledger,
		SessionID: sessionID,
	})

	enabledTools := builtinTools
	if len(opts.AllowedTools) > 0 {
		allowed := make(map[string]bool)
		for _, t := range opts.AllowedTools {
			allowed[t] = true
		}
		enabledTools = nil
		for _, t := range builtinTools {
			if allowed[t.Name] {
				enabledTools = append(enabledTools, t)
			}
		}
	}
	if len(opts.DisallowedTools) > 0 {
		denied := make(map[string]bool)
		for _, t := range opts.DisallowedTools {
			denied[t] = true
		}
		var filtered []*agentpkg.ToolDef
		for _, t := range enabledTools {
			if !denied[t.Name] {
				filtered = append(filtered, t)
			}
		}
		enabledTools = filtered
	}

	// Create agent
	agent := agentpkg.New(agentpkg.Config{
		Manifest:  manifest,
		Tools:     enabledTools,
		Breaker:   breaker,
		Ledger:    ledger,
		SessionID: sessionID,
		AgentDir:  agentDir,
	})

	// Subscribe to events and collect messages
	events := agent.Subscribe()
	done := make(chan struct{})
	go func() {
		defer close(done)
		for event := range events {
			switch event.Type {
			case agentpkg.EventDelta:
				// Skip deltas for SDK — we emit full messages at EventMessageEnd
			case agentpkg.EventMessageEnd:
				if event.Message != nil {
					result.messages = append(result.messages, Message{
						Type:    "assistant",
						Content: event.Message.Content,
						Meta: map[string]any{
							"model":       event.Message.Model,
							"provider":    event.Message.Provider,
							"stop_reason": event.Message.StopReason,
						},
					})
				}
			case agentpkg.EventToolStart:
				result.messages = append(result.messages, Message{
					Type:    "tool_use",
					Content: event.ToolName,
					Meta:    map[string]any{"args": event.ToolArgs},
				})
			case agentpkg.EventToolEnd:
				result.messages = append(result.messages, Message{
					Type:    "tool_result",
					Content: event.ToolResult,
					Meta:    map[string]any{"tool": event.ToolName, "is_error": event.IsError},
				})
			case agentpkg.EventGuardBlock:
				result.messages = append(result.messages, Message{
					Type:    "guard_block",
					Content: event.Reason,
					Meta:    map[string]any{"tool": event.ToolName},
				})
			case agentpkg.EventError:
				result.messages = append(result.messages, Message{
					Type:    "error",
					Content: event.Error.Error(),
				})
			}
		}
	}()

	// Load hooks
	hooksCfg, _ := hooks.LoadConfig(agentDir)

	// Run the agent
	result.err = agent.RunWithHooks(loopCtx, systemPrompt, opts.Prompt, hooksCfg)

	// Wait for event processing to finish
	<-done

	return result
}
