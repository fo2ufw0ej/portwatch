package anchor

import (
	"testing"
	"time"
)

func TestObserve_BelowThreshold_NotAnchored(t *testing.T) {
	s := New(3)
	for i := 0; i < 2; i++ {
		if promoted := s.Observe(8080); promoted {
			t.Fatal("expected no promotion below threshold")
		}
	}
	if s.IsAnchored(8080) {
		t.Fatal("port should not be anchored yet")
	}
}

func TestObserve_AtThreshold_Anchored(t *testing.T) {
	s := New(3)
	var promoted bool
	for i := 0; i < 3; i++ {
		promoted = s.Observe(8080)
	}
	if !promoted {
		t.Fatal("expected promotion on third observation")
	}
	if !s.IsAnchored(8080) {
		t.Fatal("port should be anchored after threshold")
	}
}

func TestObserve_BeyondThreshold_NoDoublePromotion(t *testing.T) {
	s := New(2)
	s.Observe(443)
	s.Observe(443) // anchored here
	if promoted := s.Observe(443); promoted {
		t.Fatal("should not promote again after already anchored")
	}
}

func TestReset_ClearsAnchor(t *testing.T) {
	s := New(1)
	s.Observe(22)
	if !s.IsAnchored(22) {
		t.Fatal("expected port to be anchored")
	}
	s.Reset(22)
	if s.IsAnchored(22) {
		t.Fatal("expected port to be cleared after reset")
	}
}

func TestGet_MissingPort_ReturnsFalse(t *testing.T) {
	s := New(3)
	_, ok := s.Get(9999)
	if ok {
		t.Fatal("expected missing entry")
	}
}

func TestGet_ReturnsCorrectEntry(t *testing.T) {
	s := New(5)
	before := time.Now()
	s.Observe(3000)
	e, ok := s.Get(3000)
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Streak != 1 {
		t.Fatalf("expected streak 1, got %d", e.Streak)
	}
	if e.FirstSeen.Before(before) {
		t.Fatal("FirstSeen should be >= test start")
	}
}

func TestDefaultConfig(t *testing.T) {
	c := DefaultConfig()
	if c.Threshold != 3 {
		t.Fatalf("expected default threshold 3, got %d", c.Threshold)
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("default config should be valid: %v", err)
	}
}

func TestNewFromConfig_InvalidThreshold(t *testing.T) {
	_, err := NewFromConfig(Config{Threshold: 0})
	if err == nil {
		t.Fatal("expected error for zero threshold")
	}
}

func TestNewFromConfig_Valid(t *testing.T) {
	store, err := NewFromConfig(Config{Threshold: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store == nil {
		t.Fatal("expected non-nil store")
	}
}
