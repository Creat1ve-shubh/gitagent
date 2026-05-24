package tools

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/open-gitagent/gitagent/internal/agent/util"
)

type EditTool struct {
	Root string
}

func (t *EditTool) Name() string { return "edit" }
func (t *EditTool) Description() string {
	return "Edit a file by replacing text"
}
func (t *EditTool) Schema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":        map[string]any{"type": "string"},
			"old_string":  map[string]any{"type": "string"},
			"new_string":  map[string]any{"type": "string"},
			"replace_all": map[string]any{"type": "boolean"},
		},
		"required": []string{"path", "old_string", "new_string"},
	}
}
func (t *EditTool) Metadata() Metadata { return DefaultMetadata() }
func (t *EditTool) Execute(ctx context.Context, args map[string]any) (string, error) {
	path, err := requireString(args, "path")
	if err != nil {
		return "", err
	}
	oldStr, err := requireString(args, "old_string")
	if err != nil {
		return "", err
	}
	newStr, err := requireString(args, "new_string")
	if err != nil {
		return "", err
	}
	replaceAll := false
	if v, ok := args["replace_all"].(bool); ok {
		replaceAll = v
	}

	abs, err := util.ResolvePath(t.Root, path)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(abs)
	if err != nil {
		return "", err
	}
	text := string(data)
	if oldStr == newStr {
		return "", errors.New("old_string and new_string are identical")
	}
	count := strings.Count(text, oldStr)
	if count == 0 {
		return "", errors.New("old_string not found")
	}
	if !replaceAll && count > 1 {
		return "", errors.New("old_string matches multiple times")
	}

	if replaceAll {
		text = strings.ReplaceAll(text, oldStr, newStr)
	} else {
		text = strings.Replace(text, oldStr, newStr, 1)
	}
	if err := os.WriteFile(abs, []byte(text), 0o644); err != nil {
		return "", err
	}
	return "Edited " + path, nil
}
