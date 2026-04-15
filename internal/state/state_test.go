package state_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/state"
)

func tempStorePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "state.json")
}

func TestStore_SaveAndLoad(t *testing.T) {
	store := state.NewStore(tempStorePath(t))

	snap := state.Snapshot{
		Ports:      []int{80, 443, 8080},
		RecordedAt: time.Now().UTC().Truncate(time.Second),
	}

	if err := store.Save(snap); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(loaded.Ports) != len(snap.Ports) {
		t.Fatalf("expected %d ports, got %d", len(snap.Ports), len(loaded.Ports))
	}
	for i, p := range snap.Ports {
		if loaded.Ports[i] != p {
			t.Errorf("port[%d]: want %d, got %d", i, p, loaded.Ports[i])
		}
	}
	if !loaded.RecordedAt.Equal(snap.RecordedAt) {
		t.Errorf("RecordedAt: want %v, got %v", snap.RecordedAt, loaded.RecordedAt)
	}
}

func TestStore_LoadMissing(t *testing.T) {
	store := state.NewStore(tempStorePath(t))

	snap, err := store.Load()
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(snap.Ports) != 0 {
		t.Errorf("expected empty ports, got %v", snap.Ports)
	}
}

func TestStore_OverwriteExisting(t *testing.T) {
	path := tempStorePath(t)
	store := state.NewStore(path)

	_ = store.Save(state.Snapshot{Ports: []int{22}})
	_ = store.Save(state.Snapshot{Ports: []int{22, 80}})

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded.Ports) != 2 {
		t.Errorf("expected 2 ports after overwrite, got %d", len(loaded.Ports))
	}
}

func TestStore_InvalidJSON(t *testing.T) {
	path := tempStorePath(t)
	_ = os.WriteFile(path, []byte("not-json{"), 0600)

	store := state.NewStore(path)
	_, err := store.Load()
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
