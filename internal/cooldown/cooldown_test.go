package cooldown

import (
	"testing"
	"time"
)

func TestAllow_FirstCall(t *testing.T) {
	cd := New(10 * time.Second)
	if !cd.Allow("k") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_WithinCooldown_Blocked(t *testing.T) {
	now := time.Now()
	cd := New(10 * time.Second)
	cd.now = func() time.Time { return now }
	cd.Allow("k")
	if cd.Allow("k") {
		t.Fatal("expected second call within cooldown to be blocked")
	}
}

func TestAllow_AfterCooldown_Allowed(t *testing.T) {
	now := time.Now()
	cd := New(10 * time.Second)
	cd.now = func() time.Time { return now }
	cd.Allow("k")
	cd.now = func() time.Time { return now.Add(11 * time.Second) }
	if !cd.Allow("k") {
		t.Fatal("expected call after cooldown to be allowed")
	}
}

func TestReset_ClearsEntry(t *testing.T) {
	now := time.Now()
	cd := New(10 * time.Second)
	cd.now = func() time.Time { return now }
	cd.Allow("k")
	cd.Reset("k")
	if !cd.Allow("k") {
		t.Fatal("expected allow after reset")
	}
}

func TestRemaining_WhenInCooldown(t *testing.T) {
	now := time.Now()
	cd := New(10 * time.Second)
	cd.now = func() time.Time { return now }
	cd.Allow("k")
	cd.now = func() time.Time { return now.Add(3 * time.Second) }
	r := cd.Remaining("k")
	if r <= 0 || r > 10*time.Second {
		t.Fatalf("unexpected remaining: %v", r)
	}
}

func TestRemaining_WhenNotPresent(t *testing.T) {
	cd := New(10 * time.Second)
	if cd.Remaining("missing") != 0 {
		t.Fatal("expected 0 remaining for unknown key")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Period <= 0 {
		t.Fatal("expected positive default period")
	}
}

func TestNewFromConfig_Invalid(t *testing.T) {
	_, err := NewFromConfig(Config{Period: -1})
	if err == nil {
		t.Fatal("expected error for negative period")
	}
}
