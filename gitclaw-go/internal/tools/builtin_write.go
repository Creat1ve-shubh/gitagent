package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/open-gitagent/gitclaw-go/internal/agent"
	"github.com/open-gitagent/gitclaw-go/internal/state"
)

// NewWriteTool creates the built-in write tool. Writes go through the
// ledger for conflict-free serialization and async git commits.
func NewWriteTool(cfg ToolConfig) *agent.ToolDef {
	return &agent.ToolDef{
		Name:        "write",
		Description: "Write content to a file. Creates the file and any parent directories if they don't exist. By default overwrites the file; set append: true to append instead.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "Path to the file to write (relative to agent directory)",
				},
				"content": map[string]any{
					"type":        "string",
					"description": "Content to write to the file",
				},
				"append": map[string]any{
					"type":        "boolean",
					"description": "If true, append to the file instead of overwriting",
					"default":     false,
				},
			},
			"required": []string{"path", "content"},
		},
		Execute: func(ctx context.Context, args map[string]any) (string, error) {
			path := getStringArg(args, "path")
			content := getStringArg(args, "content")
			appendMode := getBoolArg(args, "append")

			if path == "" {
				return "", fmt.Errorf("path is required")
			}

			// Resolve absolute path
			absPath := path
			if !filepath.IsAbs(path) {
				absPath = filepath.Join(cfg.AgentDir, path)
			}

			resolved, err := filepath.Abs(absPath)
			if err != nil {
				return "", fmt.Errorf("invalid path: %w", err)
			}
			agentAbs, _ := filepath.Abs(cfg.AgentDir)
			if !isSubpath(resolved, agentAbs) {
				return "", fmt.Errorf("path escapes agent directory: %s", path)
			}

			// Ensure parent directory exists
			dir := filepath.Dir(resolved)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return "", fmt.Errorf("creating directory: %w", err)
			}

			op := state.OpWrite
			if appendMode {
				op = state.OpAppend
			}

			// Use ledger for conflict-free writes
			if cfg.Ledger != nil {
				intent, err := cfg.Ledger.Acquire(resolved, op, []byte(content), "", cfg.SessionID)
				if err != nil {
					return "", fmt.Errorf("acquiring write lock: %w", err)
				}

				// Perform the actual write
				var writeErr error
				if appendMode {
					f, err := os.OpenFile(resolved, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err != nil {
						writeErr = err
					} else {
						_, writeErr = f.WriteString(content)
						f.Close()
					}
				} else {
					writeErr = os.WriteFile(resolved, []byte(content), 0644)
				}

				if writeErr != nil {
					cfg.Ledger.Fail(intent, writeErr)
					return "", fmt.Errorf("writing file: %w", writeErr)
				}

				cfg.Ledger.Complete(intent)
				return fmt.Sprintf("Wrote %d bytes to %s", len(content), path), nil
			}

			// Fallback: direct write (no ledger)
			if appendMode {
				f, err := os.OpenFile(resolved, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					return "", fmt.Errorf("opening file: %w", err)
				}
				defer f.Close()
				if _, err := f.WriteString(content); err != nil {
					return "", fmt.Errorf("appending to file: %w", err)
				}
			} else {
				if err := os.WriteFile(resolved, []byte(content), 0644); err != nil {
					return "", fmt.Errorf("writing file: %w", err)
				}
			}

			return fmt.Sprintf("Wrote %d bytes to %s", len(content), path), nil
		},
		ConcurrencySafe: false,
		ReadOnly:        false,
	}
}
