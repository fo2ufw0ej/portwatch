package shadow_test

import (
	"testing"

	"github.com/stevezaluk/portwatch/internal/shadow"
)

func TestGet_NilWhenEmpty(t *testing.T) {
	s := shadow.New()
	if s.Get() != nil {
		t.Fatal("expected nil entry on empty store")
	}
}

func TestSet_AndGet_ReturnsCopy(t *testing.T) {
	s := shadow.New()
	ports := []int{80, 443, 8080}
	s.Set(ports)

	got := s.Get()
	if got == nil {
		t.Fatal("expected non-nil entry")
	}
	if len(got.Ports) != len(ports) {
		t.Fatalf("expected %d ports, got %d", len(ports), len(got.Ports))
	}
	// Mutating original slice must not affect stored entry.
	ports[0] = 9999
	if s.Get().Ports[0] == 9999 {
		t.Fatal("store must hold an independent copy")
	}
}

func TestChanged_TrueWhenEmpty(t *testing.T) {
	s := shadow.New()
	if !s.Changed([]int{80}) {
		t.Fatal("expected Changed=true when no shadow exists")
	}
}

func TestChanged_FalseWhenSame(t *testing.T) {
	s := shadow.New()
	s.Set([]int{80, 443})
	if s.Changed([]int{443, 80}) {
		t.Fatal("expected Changed=false for same port set (different order)")
	}
}

func TestChanged_TrueWhenDifferent(t *testing.T) {
	s := shadow.New()
	s.Set([]int{80, 443})
	if !s.Changed([]int{80}) {
		t.Fatal("expected Changed=true when port removed")
	}
}

func TestChanged_TrueWhenPortAdded(t *testing.T) {
	s := shadow.New()
	s.Set([]int{80})
	if !s.Changed([]int{80, 443}) {
		t.Fatal("expected Changed=true when port added")
	}
}

func TestClear_ResetsToNil(t *testing.T) {
	s := shadow.New()
	s.Set([]int{80})
	s.Clear()
	if s.Get() != nil {
		t.Fatal("expected nil after Clear")
	}
}

func TestClear_ChangedTrueAfterClear(t *testing.T) {
	s := shadow.New()
	s.Set([]int{80})
	s.Clear()
	if !s.Changed([]int{80}) {
		t.Fatal("expected Changed=true after Clear")
	}
}
