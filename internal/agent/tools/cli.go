package tools

import (
	"context"
	"os/exec"
	"runtime"
)

type CliTool struct {
	Root string
}

func (t *CliTool) Name() string { return "cli" }
func (t *CliTool) Description() string {
	return "Execute a shell command"
}
func (t *CliTool) Schema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"command": map[string]any{"type": "string"},
		},
		"required": []string{"command"},
	}
}
func (t *CliTool) Metadata() Metadata { return DefaultMetadata() }
func (t *CliTool) Execute(ctx context.Context, args map[string]any) (string, error) {
	cmdStr, err := requireString(args, "command")
	if err != nil {
		return "", err
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "powershell", "-NoProfile", "-Command", cmdStr)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", cmdStr)
	}
	cmd.Dir = t.Root
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}
	return string(out), nil
}
