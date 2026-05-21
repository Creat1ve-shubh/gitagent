// Package config parses agent.yaml and loads the full agent configuration
// including system prompt, skills, memory, and guard settings.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/open-gitagent/gitclaw-go/internal/guard"
)

// Manifest represents the parsed agent.yaml file.
type Manifest struct {
	SpecVersion  string             `yaml:"spec_version"`
	Name         string             `yaml:"name"`
	Version      string             `yaml:"version"`
	Description  string             `yaml:"description"`
	Model        ModelConfig        `yaml:"model"`
	Tools        []string           `yaml:"tools"`
	Skills       []string           `yaml:"skills,omitempty"`
	Runtime      RuntimeConfig      `yaml:"runtime"`
	Extends      string             `yaml:"extends,omitempty"`
	Dependencies []DependencyConfig `yaml:"dependencies,omitempty"`
	Agents       map[string]SubAgentConfig `yaml:"agents,omitempty"`
	Delegation   DelegationConfig   `yaml:"delegation,omitempty"`
	Plugins      map[string]PluginConfig `yaml:"plugins,omitempty"`
	Compliance   *ComplianceConfig  `yaml:"compliance,omitempty"`
	Serve        *ServeConfig       `yaml:"serve,omitempty"`
}

// ModelConfig holds LLM model settings.
type ModelConfig struct {
	Preferred   string            `yaml:"preferred"`
	Fallback    []string          `yaml:"fallback,omitempty"`
	Constraints *ModelConstraints `yaml:"constraints,omitempty"`
}

// ModelConstraints for LLM generation parameters.
type ModelConstraints struct {
	Temperature   *float64 `yaml:"temperature,omitempty"`
	MaxTokens     *int     `yaml:"max_tokens,omitempty"`
	TopP          *float64 `yaml:"top_p,omitempty"`
	TopK          *int     `yaml:"top_k,omitempty"`
	StopSequences []string `yaml:"stop_sequences,omitempty"`
}

// RuntimeConfig holds execution limits and the guard pipeline configuration.
type RuntimeConfig struct {
	MaxTurns int          `yaml:"max_turns"`
	Timeout  int          `yaml:"timeout,omitempty"` // seconds per tool call
	Guard    *guard.Config `yaml:"guard,omitempty"`
}

// DependencyConfig for agent dependencies.
type DependencyConfig struct {
	Name    string `yaml:"name"`
	Source  string `yaml:"source"`
	Version string `yaml:"version,omitempty"`
	Mount   string `yaml:"mount,omitempty"`
}

// SubAgentConfig for delegated sub-agents.
type SubAgentConfig struct {
	Model string   `yaml:"model,omitempty"`
	Tools []string `yaml:"tools,omitempty"`
}

// DelegationConfig controls how sub-agents are dispatched.
type DelegationConfig struct {
	Mode string `yaml:"mode,omitempty"` // auto | explicit | router
}

// PluginConfig for individual plugin settings.
type PluginConfig struct {
	Enabled bool           `yaml:"enabled"`
	Config  map[string]any `yaml:"config,omitempty"`
}

// ComplianceConfig for enterprise compliance settings.
type ComplianceConfig struct {
	RiskLevel           string   `yaml:"risk_level,omitempty"`
	HumanInTheLoop      bool     `yaml:"human_in_the_loop,omitempty"`
	DataClassification  string   `yaml:"data_classification,omitempty"`
	RegulatoryFrameworks []string `yaml:"regulatory_frameworks,omitempty"`
	Recordkeeping       *RecordkeepingConfig `yaml:"recordkeeping,omitempty"`
	Review              *ReviewConfig        `yaml:"review,omitempty"`
}

// RecordkeepingConfig for audit logging settings.
type RecordkeepingConfig struct {
	AuditLogging  bool `yaml:"audit_logging,omitempty"`
	RetentionDays int  `yaml:"retention_days,omitempty"`
}

// ReviewConfig for approval workflow settings.
type ReviewConfig struct {
	RequiredApprovers int  `yaml:"required_approvers,omitempty"`
	AutoReview        bool `yaml:"auto_review,omitempty"`
}

// ServeConfig for HTTP serve mode.
type ServeConfig struct {
	Port         int               `yaml:"port,omitempty"`
	AllowedTools []string          `yaml:"allowed_tools,omitempty"`
	Constraints  *ModelConstraints `yaml:"constraints,omitempty"`
}

// LoadManifest reads and parses agent.yaml from the given directory.
func LoadManifest(agentDir string) (*Manifest, error) {
	path := filepath.Join(agentDir, "agent.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading agent.yaml: %w", err)
	}

	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parsing agent.yaml: %w", err)
	}

	// Apply defaults
	if m.Runtime.MaxTurns == 0 {
		m.Runtime.MaxTurns = 50
	}
	if m.Runtime.Timeout == 0 {
		m.Runtime.Timeout = 120
	}

	return &m, nil
}

// LoadIdentityFiles reads SOUL.md, RULES.md, and DUTIES.md and returns
// them as a combined system prompt string.
func LoadIdentityFiles(agentDir string) string {
	var parts []string

	for _, name := range []string{"SOUL.md", "RULES.md", "DUTIES.md"} {
		path := filepath.Join(agentDir, name)
		data, err := os.ReadFile(path)
		if err == nil && len(data) > 0 {
			parts = append(parts, string(data))
		}
	}

	result := ""
	for i, part := range parts {
		if i > 0 {
			result += "\n\n"
		}
		result += part
	}
	return result
}

// ParseModelString splits "provider:model-id" into provider and model ID.
// Handles the @url suffix for custom endpoints.
func ParseModelString(s string) (provider, modelID, baseURL string) {
	// Check for @url suffix: "ollama:llama3@http://localhost:11434/v1"
	atIdx := -1
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '@' {
			atIdx = i
			break
		}
	}

	modelStr := s
	if atIdx >= 0 {
		modelStr = s[:atIdx]
		baseURL = s[atIdx+1:]
	}

	// Split "provider:model-id"
	for i, c := range modelStr {
		if c == ':' {
			provider = modelStr[:i]
			modelID = modelStr[i+1:]
			return
		}
	}

	// No colon — entire string is model ID, provider is unknown
	return "", modelStr, baseURL
}
