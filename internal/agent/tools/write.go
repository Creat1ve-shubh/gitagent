package tools

import (
	"context"
	"os"
	"path/filepath"

	"github.com/open-gitagent/gitagent/internal/agent/util"
)

type WriteTool struct {
	Root string
}

func (t *WriteTool) Name() string { return "write" }
func (t *WriteTool) Description() string {
	return "Write content to a file"
}
func (t *WriteTool) Schema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":       map[string]any{"type": "string"},
			"content":    map[string]any{"type": "string"},
			"createDirs": map[string]any{"type": "boolean"},
		},
		"required": []string{"path", "content"},
	}
}
func (t *WriteTool) Metadata() Metadata { return DefaultMetadata() }
func (t *WriteTool) Execute(ctx context.Context, args map[string]any) (string, error) {
	path, err := requireString(args, "path")
	if err != nil {
		return "", err
	}
	content, err := requireString(args, "content")
	if err != nil {
		return "", err
	}
	createDirs := true
	if v, ok := args["createDirs"].(bool); ok {
		createDirs = v
	}

	abs, err := util.ResolvePath(t.Root, path)
	if err != nil {
		return "", err
	}
	if createDirs {
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			return "", err
		}
	}

	if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
		return "", err
	}
	return "Wrote " + path, nil
}
