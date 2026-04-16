package watchdog_test

import (
	"testing"
	"time"

	"github.com/danvolchek/portwatch/internal/watchdog"
)

func TestBeat_PreventsStall(t *testing.T) {
	w := watchdog.New(50 * time.Millisecond)
	w.Beat()
	if w.Stalled() {
		t.Fatal("expected not stalled immediately after Beat")
	}
}

func TestStalled_AfterTimeout(t *testing.T) {
	w := watchdog.New(20 * time.Millisecond)
	time.Sleep(40 * time.Millisecond)
	if !w.Stalled() {
		t.Fatal("expected stalled after timeout elapsed")
	}
}

func TestStalledFor_ZeroWhenHealthy(t *testing.T) {
	w := watchdog.New(100 * time.Millisecond)
	w.Beat()
	if d := w.StalledFor(); d != 0 {
		t.Fatalf("expected 0 stall duration, got %v", d)
	}
}

func TestStalledFor_PositiveWhenStalled(t *testing.T) {
	w := watchdog.New(10 * time.Millisecond)
	time.Sleep(30 * time.Millisecond)
	if d := w.StalledFor(); d <= 0 {
		t.Fatalf("expected positive stall duration, got %v", d)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := watchdog.DefaultConfig()
	if cfg.Timeout <= 0 {
		t.Fatal("default timeout must be positive")
	}
}

func TestValidate_ZeroTimeout(t *testing.T) {
	cfg := watchdog.Config{Timeout: 0}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero timeout")
	}
}

func TestNewFromConfig_Valid(t *testing.T) {
	cfg := watchdog.Config{Timeout: 30 * time.Second}
	w, err := watchdog.NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil watchdog")
	}
}

func TestNewFromConfig_Invalid(t *testing.T) {
	cfg := watchdog.Config{Timeout: -1 * time.Second}
	_, err := watchdog.NewFromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for negative timeout")
	}
}
