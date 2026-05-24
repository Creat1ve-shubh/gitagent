package guard

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Policy struct {
	AllowedTools    []string `json:"allowed_tools"`
	AllowedPaths    []string `json:"allowed_paths"`
	AllowedCommands []string `json:"allowed_commands"`
}

type Decision struct {
	Action string
	Reason string
}

func LoadPolicy(path string) (*Policy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var p Policy
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func EvaluateTool(p *Policy, tool string, args map[string]any, repoRoot string) (*Decision, error) {
	if p == nil {
		return nil, errors.New("guard policy is required")
	}
	if !contains(p.AllowedTools, tool) {
		return &Decision{Action: "block", Reason: "tool not allowlisted"}, nil
	}

	switch tool {
	case "read", "write", "edit":
		pathVal, ok := args["path"].(string)
		if !ok {
			return &Decision{Action: "block", Reason: "path missing"}, nil
		}
		abs := absPath(repoRoot, pathVal)
		if !pathAllowed(abs, repoRoot, p.AllowedPaths) {
			return &Decision{Action: "block", Reason: "path not allowlisted"}, nil
		}
	case "cli":
		cmdVal, ok := args["command"].(string)
		if !ok {
			return &Decision{Action: "block", Reason: "command missing"}, nil
		}
		if !commandAllowed(cmdVal, p.AllowedCommands) {
			return &Decision{Action: "block", Reason: "command not allowlisted"}, nil
		}
	}

	return &Decision{Action: "allow"}, nil
}

func absPath(root, p string) string {
	if filepath.IsAbs(p) {
		return filepath.Clean(p)
	}
	return filepath.Clean(filepath.Join(root, p))
}

func pathAllowed(abs string, root string, allowed []string) bool {
	root = filepath.Clean(root)
	abs = filepath.Clean(abs)
	for _, a := range allowed {
		base := a
		if strings.HasPrefix(a, "./") {
			base = filepath.Join(root, a[2:])
		} else if !filepath.IsAbs(a) {
			base = filepath.Join(root, a)
		}
		base = filepath.Clean(base)
		if abs == base || strings.HasPrefix(abs, base+string(os.PathSeparator)) {
			return true
		}
	}
	return false
}

func commandAllowed(cmd string, allowed []string) bool {
	trimmed := strings.TrimSpace(cmd)
	for _, prefix := range allowed {
		if strings.HasPrefix(trimmed, prefix) {
			return true
		}
	}
	return false
}

func contains(list []string, v string) bool {
	for _, item := range list {
		if item == v {
			return true
		}
	}
	return false
}
