package grace_test

import (
	"testing"
	"time"

	"github.com/densestvoid/portwatch/internal/grace"
)

func TestValidate_Valid(t *testing.T) {
	cfg := grace.Config{Timeout: 5 * time.Second}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_ZeroTimeout(t *testing.T) {
	cfg := grace.Config{Timeout: 0}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero timeout")
	}
}

func TestValidate_NegativeTimeout(t *testing.T) {
	cfg := grace.Config{Timeout: -1 * time.Second}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for negative timeout")
	}
}

func TestNewFromConfig_Valid(t *testing.T) {
	cfg := grace.Config{Timeout: 2 * time.Second}
	coord, err := grace.NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if coord == nil {
		t.Error("expected non-nil coordinator")
	}
}
