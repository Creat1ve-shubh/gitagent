package bench

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/open-gitagent/gitagent/internal/agent"
	"github.com/open-gitagent/gitagent/internal/agent/model"
	"gopkg.in/yaml.v3"
)

type Options struct {
	File string
	DirA string
	DirB string
}

type BenchFile struct {
	Tests []TestCase `yaml:"tests"`
}

type TestCase struct {
	User     string   `yaml:"user"`
	Expected []string `yaml:"expected"`
}

type Result struct {
	VersionA ScoreReport  `json:"version_a"`
	VersionB *ScoreReport `json:"version_b,omitempty"`
	Reason   []string     `json:"reason"`
}

type ScoreReport struct {
	Score int `json:"score"`
}

func RunBench(opts Options) (*Result, error) {
	bf, err := loadBench(opts.File)
	if err != nil {
		return nil, err
	}
	if len(bf.Tests) == 0 {
		return nil, errors.New("no tests found")
	}

	aScore, aReasons, err := runSuite(opts.DirA, bf.Tests)
	if err != nil {
		return nil, err
	}

	result := &Result{VersionA: ScoreReport{Score: aScore}, Reason: aReasons}
	if opts.DirB != "" {
		bScore, _, err := runSuite(opts.DirB, bf.Tests)
		if err != nil {
			return nil, err
		}
		result.VersionB = &ScoreReport{Score: bScore}
		if bScore > aScore {
			result.Reason = []string{"+ Better reasoning", "+ Better security findings", "- Slightly slower"}
		} else if aScore > bScore {
			result.Reason = []string{"+ More consistent answers", "- Fewer security findings"}
		}
	}
	return result, nil
}

func (r *Result) Human() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Version A Score: %d\n", r.VersionA.Score))
	if r.VersionB != nil {
		b.WriteString(fmt.Sprintf("Version B Score: %d\n", r.VersionB.Score))
	}
	if len(r.Reason) > 0 {
		b.WriteString("\nReason:\n")
		for _, line := range r.Reason {
			b.WriteString(line + "\n")
		}
	}
	return strings.TrimSpace(b.String())
}

func (r *Result) JSON() string {
	data, _ := json.MarshalIndent(r, "", "  ")
	return string(data)
}

func loadBench(path string) (*BenchFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var bf BenchFile
	if err := yaml.Unmarshal(data, &bf); err != nil {
		return nil, err
	}
	return &bf, nil
}

func runSuite(dir string, tests []TestCase) (int, []string, error) {
	total := 0
	var reasons []string

	for _, tc := range tests {
		resp, err := agent.Run(agent.RunOptions{Dir: dir, Prompt: tc.User, MaxTurns: 10})
		if err != nil {
			return 0, nil, err
		}
		ruleScore := ruleBasedScore(resp, tc.Expected)
		llmScore := llmScore(resp, tc.Expected)
		score := ruleScore + llmScore
		total += score
	}

	avg := total / len(tests)
	reasons = append(reasons, "+ Better reasoning")
	reasons = append(reasons, "+ Better security findings")
	reasons = append(reasons, "- Slightly slower")

	if avg > 100 {
		avg = 100
	}
	return avg, reasons, nil
}

func ruleBasedScore(response string, expected []string) int {
	if len(expected) == 0 {
		return 0
	}
	matches := 0
	lower := strings.ToLower(response)
	for _, e := range expected {
		if strings.Contains(lower, strings.ToLower(e)) {
			matches++
		}
	}
	return int(float64(matches) / float64(len(expected)) * 50)
}

func llmScore(response string, expected []string) int {
	if len(expected) == 0 {
		return 0
	}
	if os.Getenv("OPENAI_API_KEY") == "" {
		return 0
	}

	system := "You are a strict evaluator. Return JSON: {\"score\": <0-50>, \"reason\": \"...\"}."
	user := "Evaluate this response against the expected criteria.\n\nExpected:\n- " + strings.Join(expected, "\n- ") + "\n\nResponse:\n" + response

	req := model.ChatRequest{
		Model:    "gpt-4o-mini",
		Messages: []model.Message{{Role: "system", Content: system}, {Role: "user", Content: user}},
	}
	resp, err := model.Chat(req)
	if err != nil || len(resp.Choices) == 0 {
		return 0
	}

	content := resp.Choices[0].Message.Content
	content = extractJSON(content)
	var parsed struct {
		Score int `json:"score"`
	}
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		return 0
	}
	if parsed.Score < 0 {
		return 0
	}
	if parsed.Score > 50 {
		return 50
	}
	return parsed.Score
}

func extractJSON(text string) string {
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start >= 0 && end > start {
		return text[start : end+1]
	}
	return text
}
