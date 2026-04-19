package retrier

import (
	"testing"
	"time"
)

func TestValidate_Valid(t *testing.T) {
	cfg := DefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_ZeroAttempts(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxAttempts = 0
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero attempts")
	}
}

func TestValidate_ZeroInitialInterval(t *testing.T) {
	cfg := DefaultConfig()
	cfg.InitialInterval = 0
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero initial interval")
	}
}

func TestValidate_MaxLessThanInitial(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxInterval = cfg.InitialInterval - time.Millisecond
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error when MaxInterval < InitialInterval")
	}
}

func TestNewFromConfig_Valid(t *testing.T) {
	r, err := NewFromConfig(DefaultConfig())
	if err != nil || r == nil {
		t.Fatalf("expected valid retrier, got err=%v", err)
	}
}

func TestNewFromConfig_Invalid(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxAttempts = -1
	_, err := NewFromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}
