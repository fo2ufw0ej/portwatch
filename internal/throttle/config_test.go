package throttle

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.MinGap != 5*time.Second {
		t.Fatalf("expected 5s, got %v", cfg.MinGap)
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := Config{MinGap: time.Second}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_ZeroGap(t *testing.T) {
	cfg := Config{MinGap: 0}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero min_gap")
	}
}

func TestValidate_NegativeGap(t *testing.T) {
	cfg := Config{MinGap: -time.Second}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative min_gap")
	}
}

func TestNewFromConfig_Valid(t *testing.T) {
	cfg := Config{MinGap: 50 * time.Millisecond}
	th, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th == nil {
		t.Fatal("expected non-nil throttle")
	}
}

func TestNewFromConfig_Invalid(t *testing.T) {
	cfg := Config{MinGap: 0}
	_, err := NewFromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}
