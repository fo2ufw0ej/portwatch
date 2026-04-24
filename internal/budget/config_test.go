package budget

import (
	"testing"
	"time"
)

func TestValidate_Valid(t *testing.T) {
	cfg := Config{Window: time.Minute, Buckets: 6, Threshold: 0.99}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_ZeroWindow(t *testing.T) {
	cfg := Config{Window: 0, Buckets: 4, Threshold: 0.9}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestValidate_ZeroBuckets(t *testing.T) {
	cfg := Config{Window: time.Minute, Buckets: 0, Threshold: 0.9}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero buckets")
	}
}

func TestValidate_NegativeThreshold(t *testing.T) {
	cfg := Config{Window: time.Minute, Buckets: 4, Threshold: -0.1}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative threshold")
	}
}

func TestNewFromConfig_Valid(t *testing.T) {
	cfg := DefaultConfig()
	b, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil Budget")
	}
}

func TestNewFromConfig_Invalid(t *testing.T) {
	cfg := Config{Window: -time.Second, Buckets: 4, Threshold: 0.9}
	_, err := NewFromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}
