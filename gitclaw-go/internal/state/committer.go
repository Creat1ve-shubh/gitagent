package state

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// CommitLoop runs the serialized git commit goroutine. It drains the
// ledger's commit channel and performs git add + commit for each intent.
// This ensures git operations never contend — only one goroutine touches git.
//
// Call this in a separate goroutine: go CommitLoop(ctx, ledger)
func CommitLoop(ctx context.Context, l *Ledger) {
	for {
		select {
		case intent, ok := <-l.commitCh:
			if !ok {
				return // channel closed
			}
			if err := commitIntent(l.repoDir, intent); err != nil {
				log.Printf("[state] git commit failed for %s: %v", intent.Path, err)
				// Don't change status — it's already Committed (disk write succeeded)
				// The git commit is best-effort; the file is on disk either way.
			}

		case <-ctx.Done():
			// Drain remaining intents before exiting
			for {
				select {
				case intent, ok := <-l.commitCh:
					if !ok {
						return
					}
					_ = commitIntent(l.repoDir, intent)
				default:
					return
				}
			}
		}
	}
}

// commitIntent performs the actual git add + commit for a single intent.
func commitIntent(repoDir string, intent *WriteIntent) error {
	// Write to disk (may already be written by the tool, but we ensure it)
	switch intent.Operation {
	case OpWrite:
		dir := filepath.Dir(intent.Path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("mkdir: %w", err)
		}
		if err := os.WriteFile(intent.Path, intent.Content, 0644); err != nil {
			return fmt.Errorf("write: %w", err)
		}

	case OpAppend:
		f, err := os.OpenFile(intent.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("open for append: %w", err)
		}
		_, err = f.Write(intent.Content)
		f.Close()
		if err != nil {
			return fmt.Errorf("append: %w", err)
		}

	case OpDelete:
		if err := os.Remove(intent.Path); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("delete: %w", err)
		}
	}

	// Git add
	relPath, err := filepath.Rel(repoDir, intent.Path)
	if err != nil {
		relPath = intent.Path
	}

	addCmd := exec.Command("git", "add", relPath)
	addCmd.Dir = repoDir
	if out, err := addCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git add: %s: %w", string(out), err)
	}

	// Git commit
	msg := fmt.Sprintf("agent: %s %s", intent.Operation, filepath.Base(intent.Path))
	commitCmd := exec.Command("git", "commit", "-m", msg, "--allow-empty")
	commitCmd.Dir = repoDir
	if out, err := commitCmd.CombinedOutput(); err != nil {
		// "nothing to commit" is not an error
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return nil // nothing to commit
		}
		return fmt.Errorf("git commit: %s: %w", string(out), err)
	}

	return nil
}
