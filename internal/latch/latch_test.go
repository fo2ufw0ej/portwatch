package latch

import (
	"testing"
	"time"
)

func TestTrip_ReturnsTrueFirstTime(t *testing.T) {
	l := New()
	if !l.Trip(80) {
		t.Fatal("expected Trip to return true on first call")
	}
}

func TestTrip_ReturnsFalseIfAlreadyTripped(t *testing.T) {
	l := New()
	l.Trip(80)
	if l.Trip(80) {
		t.Fatal("expected Trip to return false when already tripped")
	}
}

func TestTripped_FalseWhenNotSet(t *testing.T) {
	l := New()
	if l.Tripped(443) {
		t.Fatal("expected Tripped to return false for unseen port")
	}
}

func TestTripped_TrueAfterTrip(t *testing.T) {
	l := New()
	l.Trip(443)
	if !l.Tripped(443) {
		t.Fatal("expected Tripped to return true after Trip")
	}
}

func TestTrippedAt_SetAfterTrip(t *testing.T) {
	l := New()
	before := time.Now()
	l.Trip(8080)
	after := time.Now()

	at, ok := l.TrippedAt(8080)
	if !ok {
		t.Fatal("expected TrippedAt ok=true")
	}
	if at.Before(before) || at.After(after) {
		t.Fatalf("TrippedAt %v outside expected range [%v, %v]", at, before, after)
	}
}

func TestTrippedAt_NotOkWhenNotSet(t *testing.T) {
	l := New()
	_, ok := l.TrippedAt(9090)
	if ok {
		t.Fatal("expected ok=false for untripped port")
	}
}

func TestReset_ClearsTrip(t *testing.T) {
	l := New()
	l.Trip(80)
	l.Reset(80)
	if l.Tripped(80) {
		t.Fatal("expected Tripped to be false after Reset")
	}
}

func TestResetAll_ClearsAll(t *testing.T) {
	l := New()
	l.Trip(80)
	l.Trip(443)
	l.Trip(8080)
	l.ResetAll()
	if l.Len() != 0 {
		t.Fatalf("expected Len 0 after ResetAll, got %d", l.Len())
	}
}

func TestLen_CountsTrippedPorts(t *testing.T) {
	l := New()
	l.Trip(80)
	l.Trip(443)
	if l.Len() != 2 {
		t.Fatalf("expected Len 2, got %d", l.Len())
	}
	l.Reset(80)
	if l.Len() != 1 {
		t.Fatalf("expected Len 1 after reset, got %d", l.Len())
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("DefaultConfig should be valid: %v", err)
	}
}

func TestNewFromConfig_Valid(t *testing.T) {
	cfg := DefaultConfig()
	l, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil Latch")
	}
}

func TestNewFromConfig_InvalidConfig(t *testing.T) {
	cfg := Config{MaxEntries: -1}
	_, err := NewFromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for negative MaxEntries")
	}
}
