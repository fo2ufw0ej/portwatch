package baseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/patrickward/portwatch/internal/baseline"
)

func tempBaselinePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "baseline.json")
}

func TestSave_AndLoad(t *testing.T) {
	store := baseline.New(tempBaselinePath(t))
	ports := []int{443, 80, 8080}

	if err := store.Save(ports); err != nil {
		t.Fatalf("Save: %v", err)
	}

	entry, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(entry.Ports) != 3 {
		t.Fatalf("expected 3 ports, got %d", len(entry.Ports))
	}
	// Ports should be sorted
	if entry.Ports[0] != 80 || entry.Ports[1] != 443 || entry.Ports[2] != 8080 {
		t.Errorf("unexpected port order: %v", entry.Ports)
	}
}

func TestLoad_Missing(t *testing.T) {
	store := baseline.New(tempBaselinePath(t))
	_, err := store.Load()
	if err == nil {
		t.Fatal("expected error for missing baseline, got nil")
	}
}

func TestSave_PreservesCreatedAt(t *testing.T) {
	path := tempBaselinePath(t)
	store := baseline.New(path)

	if err := store.Save([]int{22}); err != nil {
		t.Fatalf("first Save: %v", err)
	}
	first, _ := store.Load()

	if err := store.Save([]int{22, 443}); err != nil {
		t.Fatalf("second Save: %v", err)
	}
	second, _ := store.Load()

	if !first.CreatedAt.Equal(second.CreatedAt) {
		t.Errorf("CreatedAt changed on update: %v -> %v", first.CreatedAt, second.CreatedAt)
	}
	if !second.UpdatedAt.After(first.UpdatedAt) && !second.UpdatedAt.Equal(first.UpdatedAt) {
		t.Errorf("UpdatedAt not advanced: %v", second.UpdatedAt)
	}
}

func TestClear_RemovesFile(t *testing.T) {
	path := tempBaselinePath(t)
	store := baseline.New(path)

	_ = store.Save([]int{80})
	if err := store.Clear(); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected file to be removed after Clear")
	}
}

func TestClear_NoFile_NoError(t *testing.T) {
	store := baseline.New(tempBaselinePath(t))
	if err := store.Clear(); err != nil {
		t.Errorf("Clear on missing file should not error: %v", err)
	}
}

func TestSave_NilPorts(t *testing.T) {
	store := baseline.New(tempBaselinePath(t))
	if err := store.Save(nil); err != nil {
		t.Fatalf("Save nil: %v", err)
	}
	entry, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entry.Ports) != 0 {
		t.Errorf("expected empty ports, got %v", entry.Ports)
	}
}
