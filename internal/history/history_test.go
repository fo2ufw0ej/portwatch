package history_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/portwatch/internal/history"
)

func tempLogPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.json")
}

func TestAppend_CreatesEntry(t *testing.T) {
	log := history.New(tempLogPath(t), 0)
	if err := log.Append([]int{80, 443}, []int{22}); err != nil {
		t.Fatalf("Append failed: %v", err)
	}
	entries := log.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if len(entries[0].Opened) != 2 || len(entries[0].Closed) != 1 {
		t.Errorf("unexpected entry contents: %+v", entries[0])
	}
}

func TestAppend_PersistsToDisk(t *testing.T) {
	path := tempLogPath(t)
	log := history.New(path, 0)
	_ = log.Append([]int{8080}, nil)

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("file not written: %v", err)
	}
	var entries []history.Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("invalid JSON on disk: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 entry on disk, got %d", len(entries))
	}
}

func TestLoad_RestoresEntries(t *testing.T) {
	path := tempLogPath(t)
	log1 := history.New(path, 0)
	_ = log1.Append([]int{3000}, []int{22})
	_ = log1.Append([]int{9090}, nil)

	log2 := history.New(path, 0)
	if err := log2.Load(); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(log2.Entries()) != 2 {
		t.Errorf("expected 2 entries after load, got %d", len(log2.Entries()))
	}
}

func TestLoad_MissingFile_NoError(t *testing.T) {
	log := history.New(tempLogPath(t), 0)
	if err := log.Load(); err != nil {
		t.Errorf("expected no error for missing file, got %v", err)
	}
}

func TestMaxSize_TrimsOldEntries(t *testing.T) {
	log := history.New(tempLogPath(t), 3)
	for i := 0; i < 5; i++ {
		_ = log.Append([]int{i}, nil)
	}
	if len(log.Entries()) != 3 {
		t.Errorf("expected 3 entries (maxSize), got %d", len(log.Entries()))
	}
}
