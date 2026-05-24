package tools

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/open-gitagent/gitagent/internal/agent/util"
)

type MemoryTool struct {
	Root string
}

func (t *MemoryTool) Name() string { return "memory" }
func (t *MemoryTool) Description() string {
	return "Load or save memory from memory/MEMORY.md"
}
func (t *MemoryTool) Schema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"action":  map[string]any{"type": "string"},
			"content": map[string]any{"type": "string"},
			"message": map[string]any{"type": "string"},
		},
		"required": []string{"action"},
	}
}
func (t *MemoryTool) Metadata() Metadata { return DefaultMetadata() }
func (t *MemoryTool) Execute(ctx context.Context, args map[string]any) (string, error) {
	action, err := requireString(args, "action")
	if err != nil {
		return "", err
	}
	memPath := filepath.Join(t.Root, "memory", "MEMORY.md")

	switch action {
	case "load":
		data, err := os.ReadFile(memPath)
		if err != nil {
			return "No memories yet.", nil
		}
		if len(data) == 0 {
			return "No memories yet.", nil
		}
		return string(data), nil
	case "save":
		content, ok := args["content"].(string)
		if !ok || content == "" {
			return "", errors.New("content is required for save")
		}
		if err := util.EnsureDir(filepath.Dir(memPath)); err != nil {
			return "", err
		}
		if err := os.WriteFile(memPath, []byte(content), 0o644); err != nil {
			return "", err
		}
		return "Memory saved", nil
	default:
		return "", errors.New("unknown action")
	}
}
