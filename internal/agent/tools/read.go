package tools

import (
	"context"
	"strings"

	"github.com/open-gitagent/gitagent/internal/agent/util"
)

type ReadTool struct {
	Root string
}

func (t *ReadTool) Name() string { return "read" }
func (t *ReadTool) Description() string {
	return "Read a file with optional offset/limit"
}
func (t *ReadTool) Schema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":   map[string]any{"type": "string"},
			"offset": map[string]any{"type": "integer"},
			"limit":  map[string]any{"type": "integer"},
		},
		"required": []string{"path"},
	}
}
func (t *ReadTool) Metadata() Metadata {
	m := DefaultMetadata()
	m.ReadOnly = true
	m.ConcurrencySafe = true
	return m
}
func (t *ReadTool) Execute(ctx context.Context, args map[string]any) (string, error) {
	path, _ := args["path"].(string)
	offset := toInt(args["offset"])
	limit := toInt(args["limit"])

	abs, err := util.ResolvePath(t.Root, path)
	if err != nil {
		return "", err
	}

	lines, more, err := util.ReadLines(abs, offset, limit)
	if err != nil {
		return "", err
	}
	out := strings.Join(lines, "\n")
	if more {
		out += "\n\n[Output truncated]"
	}
	return out, nil
}
