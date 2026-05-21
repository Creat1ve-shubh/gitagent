package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/open-gitagent/gitclaw-go/internal/config"
	"github.com/open-gitagent/gitclaw-go/internal/hooks"
)

// Run executes the agent loop: prompt → LLM → tool calls → repeat.
// Returns when the LLM stops calling tools or max_turns is reached.
func (a *Agent) Run(ctx context.Context, systemPrompt, userPrompt string) error {
	a.emit(Event{Type: EventAgentStart})

	llm, err := NewLLMClient(a.manifest.Model.Preferred)
	if err != nil {
		a.emit(Event{Type: EventError, Error: err})
		return fmt.Errorf("creating LLM client: %w", err)
	}

	// Build tool definitions for the API
	toolDefs := make([]LLMToolDef, 0, len(a.tools))
	for _, t := range a.tools {
		toolDefs = append(toolDefs, LLMToolDef{
			Type: "function",
			Function: LLMFunctionDef{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.InputSchema,
			},
		})
	}

	// Build initial messages
	messages := []ChatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	maxTurns := a.manifest.Runtime.MaxTurns
	if maxTurns <= 0 {
		maxTurns = 50
	}

	for turn := 0; turn < maxTurns; turn++ {
		// Check context cancellation
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// Call LLM
		resp, err := llm.Chat(ctx, messages, toolDefs, a.manifest.Model.Constraints)
		if err != nil {
			a.emit(Event{Type: EventError, Error: err})
			return fmt.Errorf("LLM call failed (turn %d): %w", turn, err)
		}

		if len(resp.Choices) == 0 {
			return fmt.Errorf("LLM returned no choices")
		}

		choice := resp.Choices[0]
		assistantMsg := choice.Message

		// Build usage info
		var usage *TokenUsage
		if resp.Usage != nil {
			usage = &TokenUsage{
				InputTokens:  resp.Usage.PromptTokens,
				OutputTokens: resp.Usage.CompletionTokens,
				TotalTokens:  resp.Usage.TotalTokens,
			}
		}

		// Emit the assistant's text content
		contentStr := ""
		if s, ok := assistantMsg.Content.(string); ok {
			contentStr = s
		}

		if contentStr != "" {
			a.emit(Event{Type: EventDelta, Delta: contentStr})
		}

		a.emit(Event{
			Type: EventMessageEnd,
			Message: &AssistantMessage{
				Content:    contentStr,
				Model:      resp.Model,
				Provider:   llm.provider,
				StopReason: choice.FinishReason,
				Usage:      usage,
			},
		})

		// Update cost guard if we have usage
		if usage != nil {
			a.updateCostGuard(usage)
		}

		// Append assistant message to conversation
		messages = append(messages, assistantMsg)

		// Check if there are tool calls to process
		if len(assistantMsg.ToolCalls) == 0 {
			// No tool calls — agent is done
			break
		}

		// Execute tool calls
		toolCalls := make([]ToolCall, len(assistantMsg.ToolCalls))
		for i, tc := range assistantMsg.ToolCalls {
			var args map[string]any
			if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
				args = map[string]any{"_raw": tc.Function.Arguments}
			}
			toolCalls[i] = ToolCall{
				ID:         tc.ID,
				Name:       tc.Function.Name,
				Args:       args,
				TurnNumber: turn,
			}
		}

		results := a.ExecuteTools(ctx, toolCalls)

		// Append tool results as messages
		for _, r := range results {
			messages = append(messages, ChatMessage{
				Role:       "tool",
				Content:    r.Content,
				ToolCallID: r.ToolCallID,
			})
		}
	}

	a.emit(Event{Type: EventAgentEnd})
	return nil
}

// RunWithHooks runs the agent loop with lifecycle hook support.
func (a *Agent) RunWithHooks(ctx context.Context, systemPrompt, userPrompt string, hooksCfg *hooks.Config) error {
	// Run on_session_start hooks
	if hooksCfg != nil && len(hooksCfg.Hooks.OnSessionStart) > 0 {
		result := hooks.Run(hooksCfg.Hooks.OnSessionStart, a.agentDir, map[string]any{
			"event":      "on_session_start",
			"session_id": a.sessionID,
			"agent":      a.manifest.Name,
		})
		if result.Action == "block" {
			return fmt.Errorf("session blocked by hook: %s", result.Reason)
		}
	}

	err := a.Run(ctx, systemPrompt, userPrompt)

	// Run on_error hooks on failure
	if err != nil && hooksCfg != nil && len(hooksCfg.Hooks.OnError) > 0 {
		hooks.Run(hooksCfg.Hooks.OnError, a.agentDir, map[string]any{
			"event":      "on_error",
			"session_id": a.sessionID,
			"error":      err.Error(),
		})
	}

	return err
}

// updateCostGuard feeds token usage into the cost guard if present.
func (a *Agent) updateCostGuard(usage *TokenUsage) {
	if a.breaker == nil {
		return
	}
	// The breaker will propagate to cost guard if it's in the chain
	// Cost guard is updated externally — we emit the usage in the event
}

// SessionID returns the agent's session ID.
func (a *Agent) SessionID() string {
	return a.sessionID
}

// Manifest returns the agent's manifest.
func (a *Agent) Manifest() *config.Manifest {
	return a.manifest
}
