package digest_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/patrickward/portwatch/internal/digest"
)

func TestCompute_SamePortsDifferentOrder(t *testing.T) {
	a := digest.Compute([]int{80, 443, 8080})
	b := digest.Compute([]int{8080, 80, 443})
	if !digest.Equal(a, b) {
		t.Errorf("expected equal digests for same ports in different order")
	}
}

func TestCompute_DifferentPorts(t *testing.T) {
	a := digest.Compute([]int{80, 443})
	b := digest.Compute([]int{80, 8080})
	if digest.Equal(a, b) {
		t.Errorf("expected different digests for different port sets")
	}
}

func TestCompute_EmptyPorts(t *testing.T) {
	a := digest.Compute([]int{})
	b := digest.Compute([]int{})
	if !digest.Equal(a, b) {
		t.Errorf("expected equal digests for empty port sets")
	}
}

func TestChanged(t *testing.T) {
	prev := digest.Compute([]int{80})
	next := digest.Compute([]int{80, 443})
	if !digest.Changed(prev, next) {
		t.Errorf("expected Changed to return true")
	}
}

func TestPortCount(t *testing.T) {
	d := digest.Compute([]int{22, 80, 443})
	if d.PortCount != 3 {
		t.Errorf("expected PortCount 3, got %d", d.PortCount)
	}
}

func tempDigestPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "digest.json")
}

func TestStore_SaveAndLoad(t *testing.T) {
	path := tempDigestPath(t)
	s := digest.NewStore(path)
	d := digest.Compute([]int{22, 80})

	if err := s.Save(d); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !digest.Equal(d, loaded) {
		t.Errorf("loaded digest does not match saved digest")
	}
}

func TestStore_LoadMissing(t *testing.T) {
	s := digest.NewStore(filepath.Join(t.TempDir(), "missing.json"))
	_, err := s.Load()
	if err != digest.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestStore_Overwrite(t *testing.T) {
	path := tempDigestPath(t)
	s := digest.NewStore(path)

	s.Save(digest.Compute([]int{80}))
	newer := digest.Compute([]int{443})
	s.Save(newer)

	loaded, _ := s.Load()
	if !digest.Equal(newer, loaded) {
		t.Errorf("expected overwritten digest")
	}
	os.Remove(path)
}
