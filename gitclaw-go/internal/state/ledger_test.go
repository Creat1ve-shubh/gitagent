package state_test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/open-gitagent/gitclaw-go/internal/state"
)

func TestLedgerAcquireAndComplete(t *testing.T) {
	dir := t.TempDir()

	ledger := state.NewLedger(dir)
	content := []byte("hello world")

	intent, err := ledger.Acquire("test.txt", state.OpWrite, content, "tc1", "s1")
	if err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}

	if intent.ID == "" {
		t.Fatal("expected non-empty intent ID")
	}

	// Write the file
	filePath := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("writing file: %v", err)
	}

	ledger.Complete(intent)

	stats := ledger.Stats()
	if stats.Committed != 1 {
		t.Fatalf("expected 1 committed, got %d", stats.Committed)
	}
}

func TestLedgerResolveReadsFromLedger(t *testing.T) {
	dir := t.TempDir()
	ledger := state.NewLedger(dir)

	// Write version 1 to disk
	filePath := filepath.Join(dir, "data.txt")
	os.WriteFile(filePath, []byte("version1"), 0644)

	// Write version 2 through ledger
	intent, err := ledger.Acquire(filePath, state.OpWrite, []byte("version2"), "tc1", "s1")
	if err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filePath, []byte("version2"), 0644)
	ledger.Complete(intent)

	// Resolve should return version2
	data, err := ledger.Resolve(filePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "version2" {
		t.Fatalf("expected 'version2', got %q", string(data))
	}
}

func TestLedgerResolveFallsThroughToDisk(t *testing.T) {
	dir := t.TempDir()
	ledger := state.NewLedger(dir)

	// Write directly to disk (no ledger)
	filePath := filepath.Join(dir, "direct.txt")
	os.WriteFile(filePath, []byte("disk content"), 0644)

	// Resolve should read from disk
	data, err := ledger.Resolve(filePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "disk content" {
		t.Fatalf("expected 'disk content', got %q", string(data))
	}
}

func TestLedgerPerPathLocking(t *testing.T) {
	dir := t.TempDir()
	ledger := state.NewLedger(dir)

	// Two writes to DIFFERENT paths should NOT block each other
	var wg sync.WaitGroup
	results := make([]string, 2)

	wg.Add(2)
	go func() {
		defer wg.Done()
		intent, err := ledger.Acquire(filepath.Join(dir, "a.txt"), state.OpWrite, []byte("a"), "tc1", "s1")
		if err != nil {
			t.Errorf("acquire a: %v", err)
			return
		}
		os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644)
		ledger.Complete(intent)
		results[0] = "done"
	}()

	go func() {
		defer wg.Done()
		intent, err := ledger.Acquire(filepath.Join(dir, "b.txt"), state.OpWrite, []byte("b"), "tc2", "s1")
		if err != nil {
			t.Errorf("acquire b: %v", err)
			return
		}
		os.WriteFile(filepath.Join(dir, "b.txt"), []byte("b"), 0644)
		ledger.Complete(intent)
		results[1] = "done"
	}()

	wg.Wait()

	if results[0] != "done" || results[1] != "done" {
		t.Fatal("parallel writes to different paths should both complete")
	}
}

func TestLedgerSamePathSerializes(t *testing.T) {
	dir := t.TempDir()
	ledger := state.NewLedger(dir)
	filePath := filepath.Join(dir, "shared.txt")

	// Write initial content
	os.WriteFile(filePath, []byte(""), 0644)

	// Two writes to the SAME path should serialize
	order := make([]int, 0, 2)
	var mu sync.Mutex
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		intent, _ := ledger.Acquire(filePath, state.OpWrite, []byte("first"), "tc1", "s1")
		os.WriteFile(filePath, []byte("first"), 0644)
		mu.Lock()
		order = append(order, 1)
		mu.Unlock()
		ledger.Complete(intent)
	}()

	go func() {
		defer wg.Done()
		intent, _ := ledger.Acquire(filePath, state.OpWrite, []byte("second"), "tc2", "s1")
		os.WriteFile(filePath, []byte("second"), 0644)
		mu.Lock()
		order = append(order, 2)
		mu.Unlock()
		ledger.Complete(intent)
	}()

	wg.Wait()

	// Both should complete
	if len(order) != 2 {
		t.Fatalf("expected 2 completions, got %d", len(order))
	}

	stats := ledger.Stats()
	if stats.Committed != 2 {
		t.Fatalf("expected 2 committed, got %d", stats.Committed)
	}
}

func TestLedgerFailedIntentReleasesLock(t *testing.T) {
	dir := t.TempDir()
	ledger := state.NewLedger(dir)
	filePath := filepath.Join(dir, "fail.txt")

	// Acquire and fail
	intent, err := ledger.Acquire(filePath, state.OpWrite, []byte("data"), "tc1", "s1")
	if err != nil {
		t.Fatal(err)
	}
	ledger.Fail(intent, os.ErrPermission)

	// Should be able to acquire again (lock was released)
	intent2, err := ledger.Acquire(filePath, state.OpWrite, []byte("retry"), "tc2", "s1")
	if err != nil {
		t.Fatalf("second acquire should succeed after failure: %v", err)
	}
	os.WriteFile(filePath, []byte("retry"), 0644)
	ledger.Complete(intent2)

	stats := ledger.Stats()
	if stats.Failed != 1 {
		t.Fatalf("expected 1 failed, got %d", stats.Failed)
	}
	if stats.Committed != 1 {
		t.Fatalf("expected 1 committed, got %d", stats.Committed)
	}
}

func TestLedgerHasPendingWrites(t *testing.T) {
	dir := t.TempDir()
	ledger := state.NewLedger(dir)
	filePath := filepath.Join(dir, "pending.txt")

	// No pending writes initially
	if ledger.HasPendingWrites(filePath) {
		t.Fatal("should have no pending writes initially")
	}

	// Acquire creates a pending write
	intent, _ := ledger.Acquire(filePath, state.OpWrite, []byte("data"), "tc1", "s1")
	if !ledger.HasPendingWrites(filePath) {
		t.Fatal("should have pending writes after Acquire")
	}

	// Complete clears pending status
	os.WriteFile(filePath, []byte("data"), 0644)
	ledger.Complete(intent)
	if ledger.HasPendingWrites(filePath) {
		t.Fatal("should have no pending writes after Complete")
	}
}
