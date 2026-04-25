package stride

import (
	"testing"
	"time"
)

func TestRecord_IncreasesCount(t *testing.T) {
	s := New(5 * time.Second)
	s.Record()
	s.Record()
	if got := s.Count(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestRate_ZeroWhenEmpty(t *testing.T) {
	s := New(10 * time.Second)
	if r := s.Rate(); r != 0 {
		t.Fatalf("expected 0, got %f", r)
	}
}

func TestRate_CorrectValue(t *testing.T) {
	s := New(10 * time.Second)
	for i := 0; i < 5; i++ {
		s.Record()
	}
	got := s.Rate()
	expected := 5.0 / 10.0
	if got != expected {
		t.Fatalf("expected %f, got %f", expected, got)
	}
}

func TestPrune_RemovesOldEvents(t *testing.T) {
	s := New(2 * time.Second)
	now := time.Now()
	old := now.Add(-3 * time.Second)
	s.mu.Lock()
	s.times = append(s.times, old, old)
	s.mu.Unlock()
	s.Record()
	if got := s.Count(); got != 1 {
		t.Fatalf("expected 1 after pruning old events, got %d", got)
	}
}

func TestReset_ClearsEvents(t *testing.T) {
	s := New(5 * time.Second)
	s.Record()
	s.Record()
	s.Reset()
	if got := s.Count(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestDefaultConfig(t *testing.T) {
	c := DefaultConfig()
	if c.Window != 10*time.Second {
		t.Fatalf("unexpected default window: %v", c.Window)
	}
}

func TestValidate_ZeroWindow(t *testing.T) {
	c := Config{Window: 0}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewFromConfig_Valid(t *testing.T) {
	c := Config{Window: 5 * time.Second}
	s, err := NewFromConfig(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil Stride")
	}
}

func TestNewFromConfig_Invalid(t *testing.T) {
	c := Config{Window: -1 * time.Second}
	_, err := NewFromConfig(c)
	if err == nil {
		t.Fatal("expected error for negative window")
	}
}
