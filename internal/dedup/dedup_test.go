package dedup

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func diff(opened, closed []int) scanner.Diff {
	return scanner.Diff{Opened: opened, Closed: closed}
}

func TestChanged_EmptyDiff_ReturnsFalse(t *testing.T) {
	s := New()
	if s.Changed("k", diff(nil, nil)) {
		t.Fatal("expected false for empty diff")
	}
}

func TestChanged_NewDiff_ReturnsTrue(t *testing.T) {
	s := New()
	if !s.Changed("k", diff([]int{80}, nil)) {
		t.Fatal("expected true for new diff")
	}
}

func TestChanged_SameDiff_ReturnsFalse(t *testing.T) {
	s := New()
	d := diff([]int{80}, nil)
	s.Changed("k", d)
	if s.Changed("k", d) {
		t.Fatal("expected false for duplicate diff")
	}
}

func TestChanged_DifferentDiff_ReturnsTrue(t *testing.T) {
	s := New()
	s.Changed("k", diff([]int{80}, nil))
	if !s.Changed("k", diff([]int{443}, nil)) {
		t.Fatal("expected true for changed diff")
	}
}

func TestChanged_IndependentKeys(t *testing.T) {
	s := New()
	d := diff([]int{22}, nil)
	s.Changed("a", d)
	if !s.Changed("b", d) {
		t.Fatal("expected true for new key")
	}
}

func TestReset_AllowsRepeat(t *testing.T) {
	s := New()
	d := diff([]int{8080}, nil)
	s.Changed("k", d)
	s.Reset("k")
	if !s.Changed("k", d) {
		t.Fatal("expected true after reset")
	}
}

func TestDefaultConfig(t *testing.T) {
	c := DefaultConfig()
	if c.MaxKeys <= 0 {
		t.Fatalf("expected positive MaxKeys, got %d", c.MaxKeys)
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_NegativeMaxKeys(t *testing.T) {
	c := Config{MaxKeys: -1}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for negative MaxKeys")
	}
}
