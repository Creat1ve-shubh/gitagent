// Package agent implements the core agent loop — prompt → LLM → tool calls → repeat.
package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/open-gitagent/gitclaw-go/internal/config"
	"github.com/open-gitagent/gitclaw-go/internal/guard"
	"github.com/open-gitagent/gitclaw-go/internal/state"
)

// EventType identifies the kind of agent lifecycle event.
type EventType int

const (
	EventAgentStart EventType = iota
	EventDelta                // streaming text/thinking delta
	EventMessageEnd           // full assistant message assembled
	EventToolStart            // tool execution starting
	EventToolEnd              // tool execution finished
	EventAgentEnd             // agent loop finished
	EventError                // unrecoverable error
	EventGuardBlock           // tool call blocked by guard
)

// Event is emitted during the agent loop for subscribers.
type Event struct {
	Type       EventType
	Delta      string            // for EventDelta
	ToolName   string            // for EventToolStart/EventToolEnd/EventGuardBlock
	ToolArgs   map[string]any    // for EventToolStart
	ToolResult string            // for EventToolEnd
	IsError    bool              // for EventToolEnd
	Error      error             // for EventError
	Reason     string            // for EventGuardBlock
	Message    *AssistantMessage // for EventMessageEnd
}

// AssistantMessage represents a complete LLM response.
type AssistantMessage struct {
	Content    string
	Thinking   string
	Model      string
	Provider   string
	StopReason string
	Usage      *TokenUsage
}

// TokenUsage tracks token consumption.
type TokenUsage struct {
	InputTokens  int
	OutputTokens int
	TotalTokens  int
	CostUsd      float64
}

// ToolDef is a tool available to the agent.
type ToolDef struct {
	Name        string
	Description string
	InputSchema map[string]any
	Execute     func(ctx context.Context, args map[string]any) (string, error)
	// Metadata for the guard and ledger
	ConcurrencySafe bool
	ReadOnly        bool
}

// Agent orchestrates the prompt → LLM → tool → repeat loop.
type Agent struct {
	mu sync.RWMutex

	manifest  *config.Manifest
	tools     map[string]*ToolDef
	breaker   *guard.Breaker
	ledger    *state.Ledger
	sessionID string
	agentDir  string

	// subscribers receive events during the agent loop
	subscribers []chan<- Event
	subMu       sync.RWMutex
}

// Config for creating a new Agent.
type Config struct {
	Manifest  *config.Manifest
	Tools     []*ToolDef
	Breaker   *guard.Breaker
	Ledger    *state.Ledger
	SessionID string
	AgentDir  string
}

// New creates a new Agent with the given configuration.
func New(cfg Config) *Agent {
	toolMap := make(map[string]*ToolDef, len(cfg.Tools))
	for _, t := range cfg.Tools {
		toolMap[t.Name] = t
	}

	return &Agent{
		manifest:  cfg.Manifest,
		tools:     toolMap,
		breaker:   cfg.Breaker,
		ledger:    cfg.Ledger,
		sessionID: cfg.SessionID,
		agentDir:  cfg.AgentDir,
	}
}

// Subscribe returns a channel that receives agent events.
// The caller must consume from this channel to prevent blocking.
func (a *Agent) Subscribe() <-chan Event {
	ch := make(chan Event, 128)
	a.subMu.Lock()
	a.subscribers = append(a.subscribers, ch)
	a.subMu.Unlock()
	return ch
}

// emit sends an event to all subscribers (non-blocking).
func (a *Agent) emit(e Event) {
	a.subMu.RLock()
	defer a.subMu.RUnlock()
	for _, ch := range a.subscribers {
		select {
		case ch <- e:
		default:
			// subscriber channel full — drop event rather than block
		}
	}
}

// ExecuteTools runs tool calls concurrently, respecting the guard pipeline
// and write ledger. Read-only and concurrency-safe tools run in parallel;
// write tools are serialized per-path via the ledger.
func (a *Agent) ExecuteTools(ctx context.Context, calls []ToolCall) []ToolResult {
	results := make([]ToolResult, len(calls))
	var wg sync.WaitGroup

	for i, call := range calls {
		wg.Add(1)
		go func(idx int, tc ToolCall) {
			defer wg.Done()

			// 1. Guard check
			guardResp := a.breaker.Execute(guard.ToolRequest{
				ToolName:   tc.Name,
				Args:       tc.Args,
				SessionID:  a.sessionID,
				TurnNumber: tc.TurnNumber,
				Timestamp:  time.Now(),
			})

			if !guardResp.Allowed {
				a.emit(Event{
					Type:     EventGuardBlock,
					ToolName: tc.Name,
					ToolArgs: tc.Args,
					Reason:   guardResp.Reason,
				})
				results[idx] = ToolResult{
					ToolCallID: tc.ID,
					Content:    fmt.Sprintf("Tool blocked: %s", guardResp.Reason),
					IsError:    true,
				}
				return
			}

			// 2. Emit tool start
			a.emit(Event{
				Type:     EventToolStart,
				ToolName: tc.Name,
				ToolArgs: tc.Args,
			})

			// 3. Look up tool
			tool, ok := a.tools[tc.Name]
			if !ok {
				results[idx] = ToolResult{
					ToolCallID: tc.ID,
					Content:    fmt.Sprintf("Unknown tool: %s", tc.Name),
					IsError:    true,
				}
				return
			}

			// 4. Execute
			result, err := tool.Execute(ctx, tc.Args)

			// 5. Record result in breaker (for circuit breaker state)
			a.breaker.RecordResult(tc.Name, err)

			if err != nil {
				results[idx] = ToolResult{
					ToolCallID: tc.ID,
					Content:    fmt.Sprintf("Error: %v", err),
					IsError:    true,
				}
			} else {
				results[idx] = ToolResult{
					ToolCallID: tc.ID,
					Content:    result,
					IsError:    false,
				}
			}

			// 6. Emit tool end
			a.emit(Event{
				Type:       EventToolEnd,
				ToolName:   tc.Name,
				ToolResult: results[idx].Content,
				IsError:    results[idx].IsError,
			})
		}(i, call)
	}

	wg.Wait()
	return results
}

// ToolCall represents a tool invocation requested by the LLM.
type ToolCall struct {
	ID         string
	Name       string
	Args       map[string]any
	TurnNumber int
}

// ToolResult is the output of a tool execution.
type ToolResult struct {
	ToolCallID string
	Content    string
	IsError    bool
}
