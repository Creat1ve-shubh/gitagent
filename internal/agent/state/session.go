package state

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Session struct {
	Dir       string
	Branch    string
	SessionID string
}

type sessionState struct {
	SessionID string    `json:"session_id"`
	StartedAt time.Time `json:"started_at"`
}

func InitSession(dir string, session string) (*Session, error) {
	if !isGitRepo(dir) {
		return nil, errors.New("not a git repository")
	}

	if session == "" {
		session = "gitclaw/session-" + randHex(8)
	}

	if err := execGit(dir, "checkout", "-B", session); err != nil {
		return nil, err
	}

	if err := writeState(dir, session); err != nil {
		return nil, err
	}

	return &Session{Dir: dir, Branch: session, SessionID: strings.TrimPrefix(session, "gitclaw/session-")}, nil
}

func (s *Session) CommitChanges(message string) error {
	if message == "" {
		message = "gitclaw: auto-commit"
	}
	if err := execGit(s.Dir, "add", "-A"); err != nil {
		return err
	}
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	cmd.Dir = s.Dir
	if err := cmd.Run(); err == nil {
		return nil
	}
	return execGit(s.Dir, "commit", "-m", message)
}

func writeState(dir string, branch string) error {
	statePath := filepath.Join(dir, ".gitagent", "state.json")
	if err := os.MkdirAll(filepath.Dir(statePath), 0o755); err != nil {
		return err
	}
	data, _ := json.Marshal(sessionState{SessionID: branch, StartedAt: time.Now().UTC()})
	return os.WriteFile(statePath, data, 0o644)
}

func execGit(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func isGitRepo(dir string) bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = dir
	return cmd.Run() == nil
}

func randHex(n int) string {
	b := make([]byte, n/2)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
