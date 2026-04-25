package clamp

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Min != 0 {
		t.Fatalf("expected Min=0, got %d", cfg.Min)
	}
	if cfg.Max != 65535 {
		t.Fatalf("expected Max=65535, got %d", cfg.Max)
	}
}

func TestValidate_Valid(t *testing.T) {
	if err := DefaultConfig().Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_NegativeMin(t *testing.T) {
	cfg := Config{Min: -1, Max: 100}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative Min")
	}
}

func TestValidate_MaxLessThanMin(t *testing.T) {
	cfg := Config{Min: 100, Max: 50}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error when Max < Min")
	}
}

func TestNew_InvalidConfig(t *testing.T) {
	_, err := New(Config{Min: -5, Max: 10})
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}

func TestClamp_BelowMin(t *testing.T) {
	c, _ := New(Config{Min: 10, Max: 100})
	v, changed := c.Clamp(5)
	if v != 10 || !changed {
		t.Fatalf("expected (10, true), got (%d, %v)", v, changed)
	}
}

func TestClamp_AboveMax(t *testing.T) {
	c, _ := New(Config{Min: 0, Max: 1024})
	v, changed := c.Clamp(9000)
	if v != 1024 || !changed {
		t.Fatalf("expected (1024, true), got (%d, %v)", v, changed)
	}
}

func TestClamp_WithinBounds(t *testing.T) {
	c, _ := New(Config{Min: 0, Max: 65535})
	v, changed := c.Clamp(8080)
	if v != 8080 || changed {
		t.Fatalf("expected (8080, false), got (%d, %v)", v, changed)
	}
}

func TestClampAll_MixedValues(t *testing.T) {
	c, _ := New(Config{Min: 1, Max: 1000})
	out, n := c.ClampAll([]int{0, 500, 2000})
	if n != 2 {
		t.Fatalf("expected 2 changes, got %d", n)
	}
	if out[0] != 1 || out[1] != 500 || out[2] != 1000 {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestClampAll_Empty(t *testing.T) {
	c := NewDefault()
	out, n := c.ClampAll([]int{})
	if len(out) != 0 || n != 0 {
		t.Fatalf("expected empty output, got %v, %d", out, n)
	}
}

func TestBounds(t *testing.T) {
	c, _ := New(Config{Min: 5, Max: 500})
	lo, hi := c.Bounds()
	if lo != 5 || hi != 500 {
		t.Fatalf("expected (5, 500), got (%d, %d)", lo, hi)
	}
}

func TestNewFromConfig_Valid(t *testing.T) {
	_, err := NewFromConfig(DefaultConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewDefault_NoPanic(t *testing.T) {
	c := NewDefault()
	if c == nil {
		t.Fatal("expected non-nil Clamper")
	}
}
