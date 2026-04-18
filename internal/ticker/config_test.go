package ticker

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected 30s, got %v", cfg.Interval)
	}
	if cfg.Jitter != 0.1 {
		t.Errorf("expected jitter 0.1, got %v", cfg.Jitter)
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := DefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_ZeroInterval(t *testing.T) {
	cfg := Config{Interval: 0, Jitter: 0}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero interval")
	}
}

func TestValidate_NegativeInterval(t *testing.T) {
	cfg := Config{Interval: -1 * time.Second, Jitter: 0}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for negative interval")
	}
}

func TestValidate_JitterOutOfRange(t *testing.T) {
	cfg := Config{Interval: time.Second, Jitter: 1.5}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for jitter > 1")
	}
}

func TestValidate_NegativeJitter(t *testing.T) {
	cfg := Config{Interval: time.Second, Jitter: -0.1}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for negative jitter")
	}
}

func TestNew_InvalidConfig(t *testing.T) {
	_, err := New(Config{Interval: 0})
	if err == nil {
		t.Error("expected error for invalid config")
	}
}

func TestNew_ValidConfig_StopsCleanly(t *testing.T) {
	tk, err := New(Config{Interval: 100 * time.Millisecond, Jitter: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tk.Stop()
}
