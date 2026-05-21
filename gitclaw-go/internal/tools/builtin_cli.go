package tools

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/open-gitagent/gitclaw-go/internal/agent"
)

// NewCLITool creates the built-in CLI tool that runs shell commands.
func NewCLITool(cfg ToolConfig) *agent.ToolDef {
	return &agent.ToolDef{
		Name:        "cli",
		Description: "Run a shell command and return the output. Use this to execute programs, install packages, run scripts, inspect the filesystem, or perform any operation available via the command line.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"command": map[string]any{
					"type":        "string",
					"description": "The shell command to execute",
				},
			},
			"required": []string{"command"},
		},
		Execute: func(ctx context.Context, args map[string]any) (string, error) {
			command := getStringArg(args, "command")
			if command == "" {
				return "", fmt.Errorf("command is required")
			}

			timeout := cfg.Timeout
			if timeout <= 0 {
				timeout = 120
			}

			tctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
			defer cancel()

			// Use sh -c on Unix, cmd /C on Windows
			var cmd *exec.Cmd
			if runtime.GOOS == "windows" {
				cmd = exec.CommandContext(tctx, "cmd", "/C", command)
			} else {
				cmd = exec.CommandContext(tctx, "sh", "-c", command)
			}
			cmd.Dir = cfg.AgentDir

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()

			// Combine stdout and stderr
			output := stdout.String()
			if stderr.Len() > 0 {
				if output != "" {
					output += "\n"
				}
				output += stderr.String()
			}

			// Truncate very long output
			const maxOutput = 100_000
			if len(output) > maxOutput {
				output = output[:maxOutput/2] + "\n\n... (output truncated) ...\n\n" + output[len(output)-maxOutput/2:]
			}

			if err != nil {
				if tctx.Err() == context.DeadlineExceeded {
					return "", fmt.Errorf("command timed out after %ds: %s", timeout, strings.TrimSpace(output))
				}
				// Return the output even on non-zero exit
				if output != "" {
					return output, nil
				}
				return "", fmt.Errorf("command failed: %v", err)
			}

			if output == "" {
				output = "(no output)"
			}
			return output, nil
		},
		ConcurrencySafe: false,
		ReadOnly:        false,
	}
}
