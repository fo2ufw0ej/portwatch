package limiter

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.MinInterval != 5*time.Second {
		t.Fatalf("expected 5s, got %v", cfg.MinInterval)
	}
}

func TestValidate_ZeroInterval(t *testing.T) {
	cfg := Config{MinInterval: 0}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestValidate_NegativeInterval(t *testing.T) {
	cfg := Config{MinInterval: -1 * time.Second}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative interval")
	}
}

func TestNew_InvalidConfig(t *testing.T) {
	_, err := New(Config{MinInterval: 0})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAllow_FirstCall(t *testing.T) {
	l, _ := New(Config{MinInterval: 100 * time.Millisecond})
	if !l.Allow() {
		t.Fatal("first call should be allowed")
	}
}

func TestAllow_TooSoon(t *testing.T) {
	l, _ := New(Config{MinInterval: 500 * time.Millisecond})
	l.Allow()
	if l.Allow() {
		t.Fatal("second immediate call should be denied")
	}
}

func TestAllow_AfterInterval(t *testing.T) {
	l, _ := New(Config{MinInterval: 50 * time.Millisecond})
	l.Allow()
	time.Sleep(60 * time.Millisecond)
	if !l.Allow() {
		t.Fatal("call after interval should be allowed")
	}
}

func TestReset_ClearsState(t *testing.T) {
	l, _ := New(Config{MinInterval: 500 * time.Millisecond})
	l.Allow()
	l.Reset()
	if !l.Allow() {
		t.Fatal("call after reset should be allowed")
	}
}

func TestNextAllowed_BeforeFirstCall(t *testing.T) {
	l, _ := New(Config{MinInterval: 100 * time.Millisecond})
	na := l.NextAllowed()
	if na.After(time.Now().Add(time.Millisecond)) {
		t.Fatalf("expected next allowed to be now, got %v", na)
	}
}

func TestNextAllowed_AfterCall(t *testing.T) {
	l, _ := New(Config{MinInterval: 200 * time.Millisecond})
	before := time.Now()
	l.Allow()
	na := l.NextAllowed()
	if na.Before(before.Add(190 * time.Millisecond)) {
		t.Fatalf("next allowed too soon: %v", na)
	}
}
