package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/open-gitagent/gitclaw-go/internal/agent"
)

// NewReadTool creates the built-in read tool that reads file contents.
// Uses the write ledger's MVCC Resolve for read-your-own-writes consistency.
func NewReadTool(cfg ToolConfig) *agent.ToolDef {
	return &agent.ToolDef{
		Name:        "read",
		Description: "Read the contents of a file. Returns the file content as text. Supports partial reads with start/end byte offsets.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "Path to the file to read (relative to agent directory)",
				},
				"encoding": map[string]any{
					"type":        "string",
					"description": "Encoding: utf-8 (default) or base64",
					"default":     "utf-8",
				},
			},
			"required": []string{"path"},
		},
		Execute: func(ctx context.Context, args map[string]any) (string, error) {
			path := getStringArg(args, "path")
			if path == "" {
				return "", fmt.Errorf("path is required")
			}

			// Resolve absolute path
			absPath := path
			if !filepath.IsAbs(path) {
				absPath = filepath.Join(cfg.AgentDir, path)
			}

			// Security: prevent path traversal outside agent dir
			resolved, err := filepath.Abs(absPath)
			if err != nil {
				return "", fmt.Errorf("invalid path: %w", err)
			}
			agentAbs, _ := filepath.Abs(cfg.AgentDir)
			if !isSubpath(resolved, agentAbs) {
				return "", fmt.Errorf("path escapes agent directory: %s", path)
			}

			// Use ledger for MVCC reads if available
			if cfg.Ledger != nil {
				data, err := cfg.Ledger.Resolve(resolved)
				if err != nil {
					if os.IsNotExist(err) {
						return "", fmt.Errorf("file not found: %s", path)
					}
					return "", fmt.Errorf("reading file: %w", err)
				}
				return string(data), nil
			}

			// Fallback to direct read
			data, err := os.ReadFile(resolved)
			if err != nil {
				if os.IsNotExist(err) {
					return "", fmt.Errorf("file not found: %s", path)
				}
				return "", fmt.Errorf("reading file: %w", err)
			}

			return string(data), nil
		},
		ConcurrencySafe: true,
		ReadOnly:        true,
	}
}

// isSubpath checks if child is under parent directory.
func isSubpath(child, parent string) bool {
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}
	// rel must not start with ".." to be a subpath
	return rel != ".." && !filepath.IsAbs(rel) && len(rel) > 0 && rel[0] != '.'
}
