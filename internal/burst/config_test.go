package burst

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Window <= 0 {
		t.Fatalf("expected positive window, got %v", cfg.Window)
	}
	if cfg.Ceiling <= 0 {
		t.Fatalf("expected positive ceiling, got %d", cfg.Ceiling)
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := Config{Window: time.Second, Ceiling: 5}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_ZeroWindow(t *testing.T) {
	cfg := Config{Window: 0, Ceiling: 5}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestValidate_NegativeCeiling(t *testing.T) {
	cfg := Config{Window: time.Second, Ceiling: -1}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative ceiling")
	}
}

func TestNewFromConfig_Valid(t *testing.T) {
	cfg := Config{Window: time.Second, Ceiling: 3}
	tr, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil tracker")
	}
}

func TestNewFromConfig_Invalid(t *testing.T) {
	cfg := Config{Window: -time.Second, Ceiling: 3}
	_, err := NewFromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}
