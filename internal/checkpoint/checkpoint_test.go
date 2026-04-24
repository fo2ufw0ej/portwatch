package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/deanrtaylor1/portwatch/internal/checkpoint"
)

func tempCheckpointPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "checkpoint.json")
}

func TestSave_AndLoad(t *testing.T) {
	store, err := checkpoint.New(tempCheckpointPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ports := []int{22, 80, 443}
	entry, err := store.Save(ports)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	if entry.Version != 1 {
		t.Errorf("expected version 1, got %d", entry.Version)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded == nil {
		t.Fatal("expected non-nil entry")
	}
	if len(loaded.Ports) != len(ports) {
		t.Errorf("port count mismatch: got %d want %d", len(loaded.Ports), len(ports))
	}
}

func TestLoad_Missing(t *testing.T) {
	store, _ := checkpoint.New(tempCheckpointPath(t))
	entry, err := store.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry != nil {
		t.Errorf("expected nil for missing checkpoint")
	}
}

func TestSave_IncrementsVersion(t *testing.T) {
	store, _ := checkpoint.New(tempCheckpointPath(t))

	for i := uint64(1); i <= 3; i++ {
		e, err := store.Save([]int{int(i)})
		if err != nil {
			t.Fatalf("Save iteration %d: %v", i, err)
		}
		if e.Version != i {
			t.Errorf("iteration %d: expected version %d, got %d", i, i, e.Version)
		}
	}
}

func TestClear_RemovesFile(t *testing.T) {
	path := tempCheckpointPath(t)
	store, _ := checkpoint.New(path)
	_, _ = store.Save([]int{8080})

	if err := store.Clear(); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected file to be removed after Clear")
	}
}

func TestNew_EmptyPath_ReturnsError(t *testing.T) {
	_, err := checkpoint.New("")
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestConfig_DefaultIsValid(t *testing.T) {
	cfg := checkpoint.DefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("DefaultConfig should be valid, got: %v", err)
	}
}

func TestNewFromConfig_Disabled(t *testing.T) {
	cfg := checkpoint.DefaultConfig()
	cfg.Enabled = false
	store, err := checkpoint.NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store != nil {
		t.Error("expected nil store when disabled")
	}
}
