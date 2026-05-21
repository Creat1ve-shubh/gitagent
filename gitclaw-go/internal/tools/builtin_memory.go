package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/open-gitagent/gitclaw-go/internal/agent"
	"github.com/open-gitagent/gitclaw-go/internal/state"
)

// NewMemoryTool creates the built-in memory tool for persistent memory.
// Memory entries are appended to memory/MEMORY.md and git-committed.
func NewMemoryTool(cfg ToolConfig) *agent.ToolDef {
	return &agent.ToolDef{
		Name:        "memory",
		Description: "Load or save persistent memory. Use action 'load' to retrieve current memory, 'save' to append a new memory entry. Memory is git-committed for versioning.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"action": map[string]any{
					"type":        "string",
					"description": "Action: 'load' to read memory, 'save' to append entry",
					"enum":        []string{"load", "save"},
				},
				"content": map[string]any{
					"type":        "string",
					"description": "Content to save (required for 'save' action)",
				},
			},
			"required": []string{"action"},
		},
		Execute: func(ctx context.Context, args map[string]any) (string, error) {
			action := getStringArg(args, "action")
			memoryDir := filepath.Join(cfg.AgentDir, "memory")
			memoryFile := filepath.Join(memoryDir, "MEMORY.md")

			switch action {
			case "load":
				// Use ledger for MVCC reads
				if cfg.Ledger != nil {
					data, err := cfg.Ledger.Resolve(memoryFile)
					if err != nil {
						if os.IsNotExist(err) {
							return "(no memory yet)", nil
						}
						return "", fmt.Errorf("loading memory: %w", err)
					}
					content := string(data)
					if content == "" || content == "# Memory\n" {
						return "(no memory entries yet)", nil
					}
					return content, nil
				}

				data, err := os.ReadFile(memoryFile)
				if err != nil {
					if os.IsNotExist(err) {
						return "(no memory yet)", nil
					}
					return "", fmt.Errorf("loading memory: %w", err)
				}
				content := string(data)
				if content == "" || content == "# Memory\n" {
					return "(no memory entries yet)", nil
				}
				return content, nil

			case "save":
				content := getStringArg(args, "content")
				if content == "" {
					return "", fmt.Errorf("content is required for save action")
				}

				// Format the entry with timestamp
				timestamp := time.Now().Format("2006-01-02 15:04")
				entry := fmt.Sprintf("\n- [%s] %s\n", timestamp, content)

				// Ensure memory directory exists
				if err := os.MkdirAll(memoryDir, 0755); err != nil {
					return "", fmt.Errorf("creating memory directory: %w", err)
				}

				// Create file with header if it doesn't exist
				if _, err := os.Stat(memoryFile); os.IsNotExist(err) {
					if err := os.WriteFile(memoryFile, []byte("# Memory\n"), 0644); err != nil {
						return "", fmt.Errorf("creating memory file: %w", err)
					}
				}

				// Append via ledger
				if cfg.Ledger != nil {
					intent, err := cfg.Ledger.Acquire(memoryFile, state.OpAppend, []byte(entry), "", cfg.SessionID)
					if err != nil {
						return "", fmt.Errorf("acquiring memory lock: %w", err)
					}

					f, openErr := os.OpenFile(memoryFile, os.O_APPEND|os.O_WRONLY, 0644)
					if openErr != nil {
						cfg.Ledger.Fail(intent, openErr)
						return "", fmt.Errorf("opening memory file: %w", openErr)
					}
					_, writeErr := f.WriteString(entry)
					f.Close()
					if writeErr != nil {
						cfg.Ledger.Fail(intent, writeErr)
						return "", fmt.Errorf("writing memory: %w", writeErr)
					}

					cfg.Ledger.Complete(intent)
					return fmt.Sprintf("Saved to memory: %s", content), nil
				}

				// Fallback: direct append
				f, err := os.OpenFile(memoryFile, os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					return "", fmt.Errorf("opening memory file: %w", err)
				}
				defer f.Close()
				if _, err := f.WriteString(entry); err != nil {
					return "", fmt.Errorf("writing memory: %w", err)
				}

				return fmt.Sprintf("Saved to memory: %s", content), nil

			default:
				return "", fmt.Errorf("unknown action: %s (use 'load' or 'save')", action)
			}
		},
		ConcurrencySafe: false,
		ReadOnly:        false,
	}
}
