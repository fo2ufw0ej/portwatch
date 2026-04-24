package checkpoint_test

import (
	"testing"
	"time"

	"github.com/deanrtaylor1/portwatch/internal/checkpoint"
)

func TestValidate_Valid(t *testing.T) {
	cfg := checkpoint.Config{
		Path:     "/tmp/cp.json",
		Interval: time.Minute,
		Enabled:  true,
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_EmptyPath(t *testing.T) {
	cfg := checkpoint.Config{
		Path:     "",
		Interval: time.Minute,
		Enabled:  true,
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for empty path")
	}
}

func TestValidate_ZeroInterval(t *testing.T) {
	cfg := checkpoint.Config{
		Path:     "/tmp/cp.json",
		Interval: 0,
		Enabled:  true,
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero interval")
	}
}

func TestValidate_DisabledSkipsChecks(t *testing.T) {
	cfg := checkpoint.Config{
		Enabled: false,
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("disabled config should always be valid, got: %v", err)
	}
}

func TestNewFromConfig_Valid(t *testing.T) {
	dir := t.TempDir()
	cfg := checkpoint.Config{
		Path:     dir + "/cp.json",
		Interval: time.Minute,
		Enabled:  true,
	}
	store, err := checkpoint.NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store == nil {
		t.Error("expected non-nil store")
	}
}
