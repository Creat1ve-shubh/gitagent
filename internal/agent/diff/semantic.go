package diff

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type Report struct {
	BehaviorChanges []string `json:"behavior_changes"`
	RiskMedium      []string `json:"risk_medium"`
	RiskHigh        []string `json:"risk_high"`
}

func SemanticDiff(dir string) (*Report, error) {
	diff, err := gitDiff(dir)
	if err != nil {
		return nil, err
	}

	files := splitByFile(diff)
	report := &Report{}

	for file, lines := range files {
		switch category(file) {
		case "behavior":
			report.BehaviorChanges = append(report.BehaviorChanges, behaviorChanges(lines)...)
		case "permissions":
			if hasAdditions(lines) {
				report.RiskHigh = append(report.RiskHigh, "Tool execution permissions expanded")
			}
		case "memory":
			if hasAdditions(lines) || hasRemovals(lines) {
				report.RiskMedium = append(report.RiskMedium, "Memory persistence policy changed")
			}
		}
	}

	report.BehaviorChanges = unique(report.BehaviorChanges)
	report.RiskMedium = unique(report.RiskMedium)
	report.RiskHigh = unique(report.RiskHigh)

	return report, nil
}

func (r *Report) Human() string {
	var b strings.Builder
	b.WriteString("Behavior Changes\n\n")
	for _, c := range r.BehaviorChanges {
		b.WriteString("- " + c + "\n")
	}
	b.WriteString("\nRisk Assessment\n\n")
	if len(r.RiskMedium) > 0 {
		b.WriteString("Medium:\n")
		for _, c := range r.RiskMedium {
			b.WriteString("- " + c + "\n")
		}
	}
	if len(r.RiskHigh) > 0 {
		b.WriteString("\nHigh:\n")
		for _, c := range r.RiskHigh {
			b.WriteString("- " + c + "\n")
		}
	}
	return strings.TrimSpace(b.String())
}

func (r *Report) JSON() string {
	data, _ := json.MarshalIndent(r, "", "  ")
	return string(data)
}

func gitDiff(dir string) (string, error) {
	cmd := exec.Command("git", "diff", "--unified=0")
	cmd.Dir = dir
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git diff failed: %s", buf.String())
	}
	return buf.String(), nil
}

func splitByFile(diff string) map[string][]string {
	files := map[string][]string{}
	var current string
	for _, line := range strings.Split(diff, "\n") {
		if strings.HasPrefix(line, "+++ b/") {
			current = strings.TrimPrefix(line, "+++ b/")
			files[current] = []string{}
			continue
		}
		if current != "" {
			files[current] = append(files[current], line)
		}
	}
	return files
}

func category(path string) string {
	base := filepath.Base(path)
	switch {
	case base == "RULES.md" || base == "SOUL.md":
		return "behavior"
	case strings.Contains(path, "memory/"):
		return "memory"
	case base == "agent.yaml" || strings.HasSuffix(path, "guard.json") || strings.HasPrefix(path, "tools/"):
		return "permissions"
	default:
		return "other"
	}
}

func behaviorChanges(lines []string) []string {
	var out []string
	for _, l := range lines {
		if strings.HasPrefix(l, "+") && !strings.HasPrefix(l, "+++") {
			text := strings.ToLower(strings.TrimSpace(strings.TrimPrefix(l, "+")))
			if strings.Contains(text, "polite") {
				out = append(out, "Increased politeness")
			}
			if strings.Contains(text, "clarifying") || strings.Contains(text, "clarify") {
				out = append(out, "Increased willingness to ask clarifying questions")
			}
			if strings.Contains(text, "verbose") || strings.Contains(text, "detail") {
				out = append(out, "More verbose explanations")
			}
		}
		if strings.HasPrefix(l, "-") && !strings.HasPrefix(l, "---") {
			text := strings.ToLower(strings.TrimSpace(strings.TrimPrefix(l, "-")))
			if strings.Contains(text, "polite") {
				out = append(out, "Decreased politeness")
			}
		}
	}
	if len(out) == 0 {
		out = append(out, "Behavior rules updated")
	}
	return out
}

func hasAdditions(lines []string) bool {
	for _, l := range lines {
		if strings.HasPrefix(l, "+") && !strings.HasPrefix(l, "+++") {
			return true
		}
	}
	return false
}

func hasRemovals(lines []string) bool {
	for _, l := range lines {
		if strings.HasPrefix(l, "-") && !strings.HasPrefix(l, "---") {
			return true
		}
	}
	return false
}

func unique(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, v := range in {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}
