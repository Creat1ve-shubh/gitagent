package tools_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/open-gitagent/gitclaw-go/internal/state"
	"github.com/open-gitagent/gitclaw-go/internal/tools"
)

func newTestConfig(t *testing.T) tools.ToolConfig {
	t.Helper()
	dir := t.TempDir()
	ledger := state.NewLedger(dir)
	// Create memory dir
	os.MkdirAll(filepath.Join(dir, "memory"), 0755)
	os.WriteFile(filepath.Join(dir, "memory", "MEMORY.md"), []byte("# Memory\n"), 0644)

	return tools.ToolConfig{
		AgentDir:  dir,
		Timeout:   10,
		Ledger:    ledger,
		SessionID: "test-session",
	}
}

// ── CLI Tool ───────────────────────────────────────────────────────────

func TestCLIToolEcho(t *testing.T) {
	cfg := newTestConfig(t)
	cli := tools.NewCLITool(cfg)

	result, err := cli.Execute(context.Background(), map[string]any{
		"command": "echo hello",
	})
	if err != nil {
		t.Fatalf("cli tool failed: %v", err)
	}
	if !strings.Contains(result, "hello") {
		t.Fatalf("expected 'hello' in output, got %q", result)
	}
}

func TestCLIToolRequiresCommand(t *testing.T) {
	cfg := newTestConfig(t)
	cli := tools.NewCLITool(cfg)

	_, err := cli.Execute(context.Background(), map[string]any{})
	if err == nil {
		t.Fatal("expected error for missing command")
	}
}

// ── Read Tool ──────────────────────────────────────────────────────────

func TestReadTool(t *testing.T) {
	cfg := newTestConfig(t)
	read := tools.NewReadTool(cfg)

	// Create a test file
	testFile := filepath.Join(cfg.AgentDir, "test.txt")
	os.WriteFile(testFile, []byte("file content"), 0644)

	result, err := read.Execute(context.Background(), map[string]any{
		"path": "test.txt",
	})
	if err != nil {
		t.Fatalf("read tool failed: %v", err)
	}
	if result != "file content" {
		t.Fatalf("expected 'file content', got %q", result)
	}
}

func TestReadToolFileNotFound(t *testing.T) {
	cfg := newTestConfig(t)
	read := tools.NewReadTool(cfg)

	_, err := read.Execute(context.Background(), map[string]any{
		"path": "nonexistent.txt",
	})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

// ── Write Tool ─────────────────────────────────────────────────────────

func TestWriteTool(t *testing.T) {
	cfg := newTestConfig(t)
	write := tools.NewWriteTool(cfg)

	result, err := write.Execute(context.Background(), map[string]any{
		"path":    "output.txt",
		"content": "written data",
	})
	if err != nil {
		t.Fatalf("write tool failed: %v", err)
	}
	if !strings.Contains(result, "12 bytes") {
		t.Fatalf("expected byte count in result, got %q", result)
	}

	// Verify file was written
	data, err := os.ReadFile(filepath.Join(cfg.AgentDir, "output.txt"))
	if err != nil {
		t.Fatalf("reading written file: %v", err)
	}
	if string(data) != "written data" {
		t.Fatalf("expected 'written data', got %q", string(data))
	}
}

func TestWriteToolAppend(t *testing.T) {
	cfg := newTestConfig(t)
	write := tools.NewWriteTool(cfg)

	// First write
	write.Execute(context.Background(), map[string]any{
		"path":    "log.txt",
		"content": "line1\n",
	})

	// Append
	write.Execute(context.Background(), map[string]any{
		"path":    "log.txt",
		"content": "line2\n",
		"append":  true,
	})

	data, _ := os.ReadFile(filepath.Join(cfg.AgentDir, "log.txt"))
	if string(data) != "line1\nline2\n" {
		t.Fatalf("expected appended content, got %q", string(data))
	}
}

func TestWriteToolCreatesSubdirectories(t *testing.T) {
	cfg := newTestConfig(t)
	write := tools.NewWriteTool(cfg)

	_, err := write.Execute(context.Background(), map[string]any{
		"path":    "deep/nested/dir/file.txt",
		"content": "nested content",
	})
	if err != nil {
		t.Fatalf("write to nested dir failed: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(cfg.AgentDir, "deep", "nested", "dir", "file.txt"))
	if string(data) != "nested content" {
		t.Fatalf("expected 'nested content', got %q", string(data))
	}
}

// ── Memory Tool ────────────────────────────────────────────────────────

func TestMemoryToolSaveAndLoad(t *testing.T) {
	cfg := newTestConfig(t)
	mem := tools.NewMemoryTool(cfg)

	// Save
	result, err := mem.Execute(context.Background(), map[string]any{
		"action":  "save",
		"content": "user likes Go",
	})
	if err != nil {
		t.Fatalf("memory save failed: %v", err)
	}
	if !strings.Contains(result, "user likes Go") {
		t.Fatalf("expected confirmation, got %q", result)
	}

	// Load
	result, err = mem.Execute(context.Background(), map[string]any{
		"action": "load",
	})
	if err != nil {
		t.Fatalf("memory load failed: %v", err)
	}
	if !strings.Contains(result, "user likes Go") {
		t.Fatalf("expected saved memory in load result, got %q", result)
	}
}

func TestMemoryToolLoadEmpty(t *testing.T) {
	cfg := newTestConfig(t)
	mem := tools.NewMemoryTool(cfg)

	result, err := mem.Execute(context.Background(), map[string]any{
		"action": "load",
	})
	if err != nil {
		t.Fatalf("memory load failed: %v", err)
	}
	if !strings.Contains(result, "no memory") {
		t.Fatalf("expected 'no memory' message, got %q", result)
	}
}

// ── Integration: Read after Write (MVCC) ───────────────────────────────

func TestReadAfterWriteMVCC(t *testing.T) {
	cfg := newTestConfig(t)
	write := tools.NewWriteTool(cfg)
	read := tools.NewReadTool(cfg)

	// Write through ledger
	_, err := write.Execute(context.Background(), map[string]any{
		"path":    "mvcc.txt",
		"content": "ledger version",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Read should see the ledger version
	result, err := read.Execute(context.Background(), map[string]any{
		"path": "mvcc.txt",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result != "ledger version" {
		t.Fatalf("MVCC read expected 'ledger version', got %q", result)
	}
}
