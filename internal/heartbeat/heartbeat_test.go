package heartbeat

import (
	"testing"
	"time"
)

func TestNew_DefaultTTL(t *testing.T) {
	h := New(0)
	if h.ttl != 30*time.Second {
		t.Fatalf("expected default ttl 30s, got %s", h.ttl)
	}
}

func TestAlive_BeforeFirstBeat(t *testing.T) {
	h := New(5 * time.Second)
	if h.Alive() {
		t.Fatal("expected Alive=false before any beat")
	}
}

func TestAlive_AfterBeat(t *testing.T) {
	h := New(5 * time.Second)
	h.Beat()
	if !h.Alive() {
		t.Fatal("expected Alive=true immediately after beat")
	}
}

func TestAlive_AfterTTLExpires(t *testing.T) {
	now := time.Now()
	h := New(1 * time.Second)
	h.now = func() time.Time { return now }
	h.Beat()

	// Advance clock past TTL.
	h.now = func() time.Time { return now.Add(2 * time.Second) }
	if h.Alive() {
		t.Fatal("expected Alive=false after TTL expires")
	}
}

func TestLastAt_ReturnsZeroBeforeBeat(t *testing.T) {
	h := New(5 * time.Second)
	if !h.LastAt().IsZero() {
		t.Fatal("expected zero LastAt before any beat")
	}
}

func TestLastAt_ReturnsTimeOfBeat(t *testing.T) {
	fixed := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	h := New(5 * time.Second)
	h.now = func() time.Time { return fixed }
	h.Beat()
	if got := h.LastAt(); !got.Equal(fixed) {
		t.Fatalf("expected LastAt=%v, got %v", fixed, got)
	}
}

func TestStaleSince_ZeroWhenAlive(t *testing.T) {
	h := New(5 * time.Second)
	h.Beat()
	if d := h.StaleSince(); d != 0 {
		t.Fatalf("expected StaleSince=0 while alive, got %s", d)
	}
}

func TestStaleSince_PositiveWhenStale(t *testing.T) {
	now := time.Now()
	h := New(1 * time.Second)
	h.now = func() time.Time { return now }
	h.Beat()

	h.now = func() time.Time { return now.Add(3 * time.Second) }
	got := h.StaleSince()
	if got != 2*time.Second {
		t.Fatalf("expected StaleSince=2s, got %s", got)
	}
}

func TestReset_ClearsLastBeat(t *testing.T) {
	h := New(5 * time.Second)
	h.Beat()
	h.Reset()
	if h.Alive() {
		t.Fatal("expected Alive=false after Reset")
	}
	if !h.LastAt().IsZero() {
		t.Fatal("expected zero LastAt after Reset")
	}
}
