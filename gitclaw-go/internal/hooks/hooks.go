// Package hooks implements the lifecycle hook system — shell scripts and
// programmatic functions that intercept agent events.
package hooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// HookDefinition is a single hook entry from hooks.yaml.
type HookDefinition struct {
	Script      string `yaml:"script" json:"script"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
}

// Config holds all hooks loaded from hooks/hooks.yaml.
type Config struct {
	Hooks struct {
		OnSessionStart  []HookDefinition `yaml:"on_session_start,omitempty"`
		PreToolUse      []HookDefinition `yaml:"pre_tool_use,omitempty"`
		PostToolFailure []HookDefinition `yaml:"post_tool_failure,omitempty"`
		PostResponse    []HookDefinition `yaml:"post_response,omitempty"`
		PreQuery        []HookDefinition `yaml:"pre_query,omitempty"`
		FileChanged     []HookDefinition `yaml:"file_changed,omitempty"`
		OnError         []HookDefinition `yaml:"on_error,omitempty"`
	} `yaml:"hooks"`
}

// Result is the outcome of running a hook chain.
type Result struct {
	Action string         `json:"action"` // "allow", "block", "modify"
	Reason string         `json:"reason,omitempty"`
	Args   map[string]any `json:"args,omitempty"`
}

// LoadConfig reads hooks/hooks.yaml from the agent directory.
func LoadConfig(agentDir string) (*Config, error) {
	path := filepath.Join(agentDir, "hooks", "hooks.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // no hooks configured
		}
		return nil, fmt.Errorf("reading hooks config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing hooks config: %w", err)
	}

	return &cfg, nil
}

// Run executes a chain of hook definitions with the given input.
// Returns on the first "block" or "modify" result. If all hooks pass, returns "allow".
func Run(hooks []HookDefinition, agentDir string, input map[string]any) Result {
	if len(hooks) == 0 {
		return Result{Action: "allow"}
	}

	for _, hook := range hooks {
		result, err := executeHook(hook, agentDir, input)
		if err != nil {
			// Hook errors don't block by default — log and continue
			fmt.Fprintf(os.Stderr, "[hooks] error in %q: %v\n", hook.Script, err)
			continue
		}
		if result.Action == "block" || result.Action == "modify" {
			return result
		}
	}

	return Result{Action: "allow"}
}

// executeHook runs a single hook script, passing JSON on stdin and reading
// JSON from stdout. Times out after 10 seconds.
func executeHook(hook HookDefinition, agentDir string, input map[string]any) (Result, error) {
	scriptPath := filepath.Join(agentDir, "hooks", hook.Script)

	// Path traversal guard
	resolved, err := filepath.Abs(scriptPath)
	if err != nil {
		return Result{}, fmt.Errorf("resolving script path: %w", err)
	}
	base, _ := filepath.Abs(agentDir)
	rel, err := filepath.Rel(base, resolved)
	if err != nil || len(rel) < 1 || rel[0] == '.' {
		return Result{}, fmt.Errorf("hook script %q escapes agent directory", hook.Script)
	}

	// Check script exists
	if _, err := os.Stat(resolved); err != nil {
		return Result{}, fmt.Errorf("hook script not found: %s", hook.Script)
	}

	// Marshal input
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return Result{}, fmt.Errorf("marshaling hook input: %w", err)
	}

	// Execute with timeout
	cmd := exec.Command("sh", resolved)
	cmd.Dir = agentDir
	cmd.Stdin = bytes.NewReader(inputJSON)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Start with timeout
	done := make(chan error, 1)
	if err := cmd.Start(); err != nil {
		return Result{}, fmt.Errorf("starting hook: %w", err)
	}
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		if err != nil {
			return Result{}, fmt.Errorf("hook %q failed (exit %v): %s", hook.Script, err, stderr.String())
		}
	case <-time.After(10 * time.Second):
		cmd.Process.Kill()
		return Result{}, fmt.Errorf("hook %q timed out after 10s", hook.Script)
	}

	// Parse output
	output := bytes.TrimSpace(stdout.Bytes())
	if len(output) == 0 {
		return Result{Action: "allow"}, nil
	}

	var result Result
	if err := json.Unmarshal(output, &result); err != nil {
		// Non-JSON output → treat as allow
		return Result{Action: "allow"}, nil
	}

	return result, nil
}
