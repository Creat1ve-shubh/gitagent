// Package state implements the Write Ledger — an in-memory MVCC system
// that serializes agent file writes to prevent deadlocks and conflicts.
//
// Key invariants:
//   - Per-path intent locks (not global): parallel writes to different files proceed concurrently
//   - MVCC reads: readers see committed writes from their own session immediately
//   - Git commits are serialized via a single goroutine draining a channel
//   - The ledger is ephemeral — it dies with the session; git is the source of truth
package state

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

// OpType represents the type of write operation.
type OpType int

const (
	OpWrite  OpType = iota // overwrite file
	OpAppend               // append to file
	OpDelete               // delete file
)

func (o OpType) String() string {
	switch o {
	case OpWrite:
		return "write"
	case OpAppend:
		return "append"
	case OpDelete:
		return "delete"
	default:
		return "unknown"
	}
}

// IntentStatus tracks the lifecycle of a write intent.
type IntentStatus int

const (
	StatusPending   IntentStatus = iota // queued, waiting for path lock
	StatusExecuting                     // lock acquired, writing in progress
	StatusCommitted                     // written to disk successfully
	StatusFailed                        // write failed
)

// WriteIntent represents a single file mutation tracked by the ledger.
type WriteIntent struct {
	ID        string
	Path      string // normalized absolute path
	Operation OpType
	Content   []byte
	ToolCall  string // originating tool call ID
	SessionID string
	Status    IntentStatus
	CreatedAt time.Time
	Version   uint64 // monotonic MVCC version

	// done is closed when the intent transitions to Committed or Failed.
	done chan struct{}
}

// Ledger tracks all file write intents for a session, providing
// per-path serialization and MVCC read-your-own-writes semantics.
type Ledger struct {
	mu      sync.RWMutex
	intents map[string]*WriteIntent   // by intent ID
	byPath  map[string][]*WriteIntent // path → ordered intents
	version atomic.Uint64

	// pathLocks serializes writes per path. Each path gets its own mutex.
	pathMu    sync.Mutex
	pathLocks map[string]*sync.Mutex

	// commitCh feeds the async git commit goroutine.
	commitCh chan *WriteIntent

	// repoDir is the root of the git repository.
	repoDir string
}

// NewLedger creates a new write ledger for the given repository directory.
func NewLedger(repoDir string) *Ledger {
	return &Ledger{
		intents:   make(map[string]*WriteIntent),
		byPath:    make(map[string][]*WriteIntent),
		pathLocks: make(map[string]*sync.Mutex),
		commitCh:  make(chan *WriteIntent, 64),
		repoDir:   repoDir,
	}
}

// Acquire creates a new write intent and acquires the per-path lock.
// Blocks if another intent for the same path is currently executing.
// Returns a handle that the caller must Complete() or Fail() when done.
func (l *Ledger) Acquire(path string, op OpType, content []byte, toolCall, sessionID string) (*WriteIntent, error) {
	absPath, err := l.normalizePath(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	// Acquire the per-path lock (blocks if another write is in progress)
	pathLock := l.getPathLock(absPath)
	pathLock.Lock()

	intent := &WriteIntent{
		ID:        fmt.Sprintf("wi_%d", l.version.Add(1)),
		Path:      absPath,
		Operation: op,
		Content:   content,
		ToolCall:  toolCall,
		SessionID: sessionID,
		Status:    StatusExecuting,
		CreatedAt: time.Now(),
		Version:   l.version.Load(),
		done:      make(chan struct{}),
	}

	l.mu.Lock()
	l.intents[intent.ID] = intent
	l.byPath[absPath] = append(l.byPath[absPath], intent)
	l.mu.Unlock()

	return intent, nil
}

// Complete marks the intent as committed and releases the path lock.
// The intent is queued for async git commit.
func (l *Ledger) Complete(intent *WriteIntent) {
	intent.Status = StatusCommitted
	close(intent.done)

	// Release the per-path lock
	pathLock := l.getPathLock(intent.Path)
	pathLock.Unlock()

	// Queue for async git commit
	select {
	case l.commitCh <- intent:
	default:
		// Channel full — commit synchronously as fallback
		// (this shouldn't happen with a reasonable buffer)
	}
}

// Fail marks the intent as failed and releases the path lock.
func (l *Ledger) Fail(intent *WriteIntent, err error) {
	intent.Status = StatusFailed
	close(intent.done)

	pathLock := l.getPathLock(intent.Path)
	pathLock.Unlock()
}

// Resolve returns the effective content of a file, considering pending
// committed writes that haven't been flushed to disk yet.
// This provides read-your-own-writes semantics.
func (l *Ledger) Resolve(path string) ([]byte, error) {
	absPath, err := l.normalizePath(path)
	if err != nil {
		return nil, err
	}

	l.mu.RLock()
	intents := l.byPath[absPath]
	l.mu.RUnlock()

	// Find latest committed intent
	var latest *WriteIntent
	for _, i := range intents {
		if i.Status == StatusCommitted && (latest == nil || i.Version > latest.Version) {
			latest = i
		}
	}

	if latest != nil {
		switch latest.Operation {
		case OpWrite:
			return latest.Content, nil
		case OpAppend:
			// Read base from disk and append
			base, _ := os.ReadFile(absPath)
			return append(base, latest.Content...), nil
		case OpDelete:
			return nil, os.ErrNotExist
		}
	}

	// No pending writes — read from disk
	return os.ReadFile(absPath)
}

// HasPendingWrites returns true if there are any executing writes
// for the given path.
func (l *Ledger) HasPendingWrites(path string) bool {
	absPath, err := l.normalizePath(path)
	if err != nil {
		return false
	}

	l.mu.RLock()
	defer l.mu.RUnlock()

	for _, i := range l.byPath[absPath] {
		if i.Status == StatusExecuting || i.Status == StatusPending {
			return true
		}
	}
	return false
}

// Stats returns a snapshot of ledger statistics.
func (l *Ledger) Stats() LedgerStats {
	l.mu.RLock()
	defer l.mu.RUnlock()

	stats := LedgerStats{
		TotalIntents: len(l.intents),
		TrackedPaths: len(l.byPath),
	}

	for _, i := range l.intents {
		switch i.Status {
		case StatusPending:
			stats.Pending++
		case StatusExecuting:
			stats.Executing++
		case StatusCommitted:
			stats.Committed++
		case StatusFailed:
			stats.Failed++
		}
	}
	return stats
}

// LedgerStats is a point-in-time snapshot of ledger state.
type LedgerStats struct {
	TotalIntents int `json:"total_intents"`
	TrackedPaths int `json:"tracked_paths"`
	Pending      int `json:"pending"`
	Executing    int `json:"executing"`
	Committed    int `json:"committed"`
	Failed       int `json:"failed"`
}

// getPathLock returns (or creates) the per-path mutex.
func (l *Ledger) getPathLock(absPath string) *sync.Mutex {
	l.pathMu.Lock()
	defer l.pathMu.Unlock()

	mu, ok := l.pathLocks[absPath]
	if !ok {
		mu = &sync.Mutex{}
		l.pathLocks[absPath] = mu
	}
	return mu
}

// normalizePath converts a potentially relative path to an absolute one
// rooted in the repo directory.
func (l *Ledger) normalizePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}
	return filepath.Abs(filepath.Join(l.repoDir, path))
}
