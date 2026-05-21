package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-gitagent/gitclaw-go/internal/config"
)

func TestLoadManifest(t *testing.T) {
	dir := t.TempDir()
	yamlContent := `
spec_version: "0.1.0"
name: test-agent
version: 1.0.0
description: A test agent
model:
  preferred: "openai:gpt-4o"
  fallback:
    - "groq:llama-3.3-70b"
tools:
  - cli
  - read
  - write
  - memory
runtime:
  max_turns: 30
  timeout: 60
`
	os.WriteFile(filepath.Join(dir, "agent.yaml"), []byte(yamlContent), 0644)

	m, err := config.LoadManifest(dir)
	if err != nil {
		t.Fatalf("LoadManifest: %v", err)
	}

	if m.Name != "test-agent" {
		t.Fatalf("expected name 'test-agent', got %q", m.Name)
	}
	if m.Model.Preferred != "openai:gpt-4o" {
		t.Fatalf("expected model 'openai:gpt-4o', got %q", m.Model.Preferred)
	}
	if len(m.Tools) != 4 {
		t.Fatalf("expected 4 tools, got %d", len(m.Tools))
	}
	if m.Runtime.MaxTurns != 30 {
		t.Fatalf("expected max_turns 30, got %d", m.Runtime.MaxTurns)
	}
}

func TestLoadManifestDefaults(t *testing.T) {
	dir := t.TempDir()
	yamlContent := `
name: minimal
model:
  preferred: "openai:gpt-4o-mini"
`
	os.WriteFile(filepath.Join(dir, "agent.yaml"), []byte(yamlContent), 0644)

	m, err := config.LoadManifest(dir)
	if err != nil {
		t.Fatalf("LoadManifest: %v", err)
	}

	if m.Runtime.MaxTurns != 50 {
		t.Fatalf("expected default max_turns 50, got %d", m.Runtime.MaxTurns)
	}
	if m.Runtime.Timeout != 120 {
		t.Fatalf("expected default timeout 120, got %d", m.Runtime.Timeout)
	}
}

func TestLoadManifestWithGuard(t *testing.T) {
	dir := t.TempDir()
	yamlContent := `
name: guarded
model:
  preferred: "openai:gpt-4o"
runtime:
  max_turns: 10
  guard:
    rate_limit:
      max_per_window: 50
    circuit_breaker:
      failure_threshold: 3
    policies:
      - tool: cli
        deny: "args.command matches 'rm -rf.*'"
`
	os.WriteFile(filepath.Join(dir, "agent.yaml"), []byte(yamlContent), 0644)

	m, err := config.LoadManifest(dir)
	if err != nil {
		t.Fatalf("LoadManifest: %v", err)
	}

	if m.Runtime.Guard == nil {
		t.Fatal("expected guard config to be loaded")
	}
	if len(m.Runtime.Guard.Policies) != 1 {
		t.Fatalf("expected 1 policy, got %d", len(m.Runtime.Guard.Policies))
	}
}

func TestParseModelString(t *testing.T) {
	tests := []struct {
		input    string
		provider string
		modelID  string
		baseURL  string
	}{
		{"openai:gpt-4o", "openai", "gpt-4o", ""},
		{"anthropic:claude-sonnet-4-6", "anthropic", "claude-sonnet-4-6", ""},
		{"ollama:llama3@http://localhost:11434/v1", "ollama", "llama3", "http://localhost:11434/v1"},
		{"gpt-4o-mini", "", "gpt-4o-mini", ""},
	}

	for _, tt := range tests {
		provider, modelID, baseURL := config.ParseModelString(tt.input)
		if provider != tt.provider {
			t.Errorf("ParseModelString(%q): provider = %q, want %q", tt.input, provider, tt.provider)
		}
		if modelID != tt.modelID {
			t.Errorf("ParseModelString(%q): modelID = %q, want %q", tt.input, modelID, tt.modelID)
		}
		if baseURL != tt.baseURL {
			t.Errorf("ParseModelString(%q): baseURL = %q, want %q", tt.input, baseURL, tt.baseURL)
		}
	}
}

func TestLoadIdentityFiles(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "SOUL.md"), []byte("# Soul\nI am an agent"), 0644)
	os.WriteFile(filepath.Join(dir, "RULES.md"), []byte("# Rules\nBe safe"), 0644)

	prompt := config.LoadIdentityFiles(dir)
	if prompt == "" {
		t.Fatal("expected non-empty system prompt")
	}
	if len(prompt) < 10 {
		t.Fatalf("system prompt too short: %q", prompt)
	}
}
