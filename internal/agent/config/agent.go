package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Manifest struct {
	SpecVersion string `yaml:"spec_version"`
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
	Model       struct {
		Preferred   string   `yaml:"preferred"`
		Fallback    []string `yaml:"fallback"`
		Constraints struct {
			Temperature   *float64 `yaml:"temperature"`
			MaxTokens     *int     `yaml:"max_tokens"`
			TopP          *float64 `yaml:"top_p"`
			TopK          *int     `yaml:"top_k"`
			StopSequences []string `yaml:"stop_sequences"`
		} `yaml:"constraints"`
	} `yaml:"model"`
	Tools   []string `yaml:"tools"`
	Runtime struct {
		MaxTurns int `yaml:"max_turns"`
		Timeout  int `yaml:"timeout"`
	} `yaml:"runtime"`
}

func LoadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}
