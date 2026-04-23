package probe_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/probe"
)

func TestValidate_Valid(t *testing.T) {
	cfg := probe.DefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_EmptyHost(t *testing.T) {
	cfg := probe.DefaultConfig()
	cfg.Host = ""
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty host")
	}
}

func TestValidate_ZeroTimeout(t *testing.T) {
	cfg := probe.DefaultConfig()
	cfg.Timeout = 0
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero timeout")
	}
}

func TestValidate_NegativeTimeout(t *testing.T) {
	cfg := probe.DefaultConfig()
	cfg.Timeout = -time.Second
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative timeout")
	}
}

func TestValidate_ZeroMaxFailures(t *testing.T) {
	cfg := probe.DefaultConfig()
	cfg.MaxFailures = 0
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero max_failures")
	}
}

func TestNewFromConfig_Valid(t *testing.T) {
	p, err := probe.NewFromConfig(probe.DefaultConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}

func TestNewFromConfig_Invalid(t *testing.T) {
	cfg := probe.DefaultConfig()
	cfg.MaxFailures = -1
	_, err := probe.NewFromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}
