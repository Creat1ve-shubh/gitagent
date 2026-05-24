package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/open-gitagent/gitagent/internal/agent/config"
	"github.com/open-gitagent/gitagent/internal/agent/guard"
	"github.com/open-gitagent/gitagent/internal/agent/model"
	"github.com/open-gitagent/gitagent/internal/agent/state"
	"github.com/open-gitagent/gitagent/internal/agent/tools"
)

type RunOptions struct {
	Dir      string
	Prompt   string
	Model    string
	MaxTurns int
}

type ToolRegistry struct {
	Tools []tools.Tool
}

func Run(opts RunOptions) (string, error) {
	manifest, err := config.LoadManifest(filepath.Join(opts.Dir, "agent.yaml"))
	if err != nil {
		return "", err
	}

	if opts.MaxTurns <= 0 {
		if manifest.Runtime.MaxTurns > 0 {
			opts.MaxTurns = manifest.Runtime.MaxTurns
		} else {
			opts.MaxTurns = 30
		}
	}

	session, err := state.InitSession(opts.Dir, "")
	if err != nil {
		return "", err
	}

	guardPolicy, err := guard.LoadPolicy(filepath.Join(opts.Dir, "agents", manifest.Name, "guard.json"))
	if err != nil {
		return "", err
	}

	registry := ToolRegistry{Tools: builtinTools(opts.Dir, manifest.Tools)}

	systemPrompt := buildSystemPrompt(opts.Dir)
	messages := []model.Message{{Role: "system", Content: systemPrompt}, {Role: "user", Content: opts.Prompt}}

	modelName := manifest.Model.Preferred
	if opts.Model != "" {
		modelName = opts.Model
	}

	lock := state.NewLock(opts.Dir, 60*time.Second)
	ctx := context.Background()

	for i := 0; i < opts.MaxTurns; i++ {
		toolDefs := toolDefsFromRegistry(registry)
		req := model.ChatRequest{Model: modelName, Messages: messages, Tools: toolDefs}
		resp, err := model.Chat(req)
		if err != nil {
			return "", err
		}
		if len(resp.Choices) == 0 {
			return "", errors.New("no response from model")
		}

		choice := resp.Choices[0].Message
		if len(choice.ToolCalls) == 0 {
			return choice.Content, nil
		}

		messages = append(messages, model.Message{Role: "assistant", Content: choice.Content})

		results := make([]string, len(choice.ToolCalls))
		errCh := make(chan error, len(choice.ToolCalls))
		allReadOnly := true
		for _, call := range choice.ToolCalls {
			tool := findTool(registry, call.Function.Name)
			if tool == nil {
				return "", fmt.Errorf("unknown tool: %s", call.Function.Name)
			}
			meta := tool.Metadata()
			if !meta.ReadOnly || !meta.ConcurrencySafe {
				allReadOnly = false
				break
			}
		}

		if allReadOnly {
			for i, call := range choice.ToolCalls {
				call := call
				idx := i
				go func() {
					res, err := executeTool(ctx, call, registry, guardPolicy, opts.Dir, lock, session)
					if err != nil {
						errCh <- err
						return
					}
					results[idx] = res
					errCh <- nil
				}()
			}
			for range choice.ToolCalls {
				if err := <-errCh; err != nil {
					return "", err
				}
			}
		} else {
			for i, call := range choice.ToolCalls {
				res, err := executeTool(ctx, call, registry, guardPolicy, opts.Dir, lock, session)
				if err != nil {
					return "", err
				}
				results[i] = res
			}
		}

		for i, call := range choice.ToolCalls {
			messages = append(messages, model.Message{Role: "tool", ToolID: call.ID, Name: call.Function.Name, Content: results[i]})
		}
	}

	return "", errors.New("max turns exceeded")
}

func builtinTools(root string, allow []string) []tools.Tool {
	all := []tools.Tool{
		&tools.ReadTool{Root: root},
		&tools.WriteTool{Root: root},
		&tools.EditTool{Root: root},
		&tools.MemoryTool{Root: root},
		&tools.CliTool{Root: root},
	}
	if len(allow) == 0 {
		return all
	}
	var filtered []tools.Tool
	for _, t := range all {
		for _, name := range allow {
			if t.Name() == name {
				filtered = append(filtered, t)
			}
		}
	}
	return filtered
}

func buildSystemPrompt(root string) string {
	parts := []string{}
	files := []string{"SOUL.md", "RULES.md", filepath.Join("memory", "MEMORY.md")}
	for _, f := range files {
		data, err := os.ReadFile(filepath.Join(root, f))
		if err == nil {
			parts = append(parts, string(data))
		}
	}
	return strings.Join(parts, "\n\n")
}

func toolDefsFromRegistry(reg ToolRegistry) []model.ToolDef {
	var defs []model.ToolDef
	for _, t := range reg.Tools {
		defs = append(defs, model.ToolDef{
			Type: "function",
			Function: map[string]any{
				"name":        t.Name(),
				"description": t.Description(),
				"parameters":  t.Schema(),
			},
		})
	}
	return defs
}

func findTool(reg ToolRegistry, name string) tools.Tool {
	for _, t := range reg.Tools {
		if t.Name() == name {
			return t
		}
	}
	return nil
}

func executeTool(
	ctx context.Context,
	call model.ToolCall,
	registry ToolRegistry,
	guardPolicy *guard.Policy,
	repoRoot string,
	lock *state.Lock,
	session *state.Session,
) (string, error) {
	args := map[string]any{}
	if err := json.Unmarshal([]byte(call.Function.Arguments), &args); err != nil {
		return "", err
	}

	decision, err := guard.EvaluateTool(guardPolicy, call.Function.Name, args, repoRoot)
	if err != nil {
		return "", err
	}
	if decision.Action != "allow" {
		return "", fmt.Errorf("tool blocked: %s", decision.Reason)
	}

	tool := findTool(registry, call.Function.Name)
	if tool == nil {
		return "", fmt.Errorf("unknown tool: %s", call.Function.Name)
	}

	meta := tool.Metadata()
	if !meta.ReadOnly {
		if err := lock.Acquire(); err != nil {
			return "", err
		}
	}

	result, err := tool.Execute(ctx, args)
	if !meta.ReadOnly {
		_ = lock.Release()
		_ = session.CommitChanges("gitclaw: tool " + tool.Name())
	}
	if err != nil {
		return "", err
	}
	if meta.MaxResultChars > 0 && len(result) > meta.MaxResultChars {
		result = result[:meta.MaxResultChars] + "\n\n[Truncated]"
	}

	return result, nil
}
