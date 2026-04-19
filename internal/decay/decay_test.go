package decay

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.HalfLife <= 0 {
		t.Fatal("expected positive HalfLife")
	}
}

func TestValidate_InvalidHalfLife(t *testing.T) {
	cfg := Config{HalfLife: 0}
	if err := cfg.Validate(); err != ErrInvalidHalfLife {
		t.Fatalf("expected ErrInvalidHalfLife, got %v", err)
	}
}

func TestNew_InvalidConfig(t *testing.T) {
	_, err := New(Config{HalfLife: -1})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAdd_AndGet_ImmediatelyFresh(t *testing.T) {
	s, err := New(Config{HalfLife: time.Hour})
	if err != nil {
		t.Fatal(err)
	}
	s.Add("port:8080", 1.0)
	v := s.Get("port:8080")
	if v < 0.99 || v > 1.01 {
		t.Fatalf("expected ~1.0, got %f", v)
	}
}

func TestGet_Missing_ReturnsZero(t *testing.T) {
	s, _ := NewDefault()
	if v := s.Get("unknown"); v != 0 {
		t.Fatalf("expected 0, got %f", v)
	}
}

func TestDecay_ReducesOverTime(t *testing.T) {
	cfg := Config{HalfLife: 100 * time.Millisecond}
	s, _ := New(cfg)
	s.Add("k", 1.0)
	time.Sleep(200 * time.Millisecond)
	v := s.Get("k")
	if v >= 0.5 {
		t.Fatalf("expected value < 0.5 after one half-life, got %f", v)
	}
}

func TestReset_ClearsScore(t *testing.T) {
	s, _ := NewDefault()
	s.Add("k", 5.0)
	s.Reset("k")
	if v := s.Get("k"); v != 0 {
		t.Fatalf("expected 0 after reset, got %f", v)
	}
}

func TestAdd_Accumulates(t *testing.T) {
	s, _ := New(Config{HalfLife: time.Hour})
	s.Add("k", 1.0)
	s.Add("k", 1.0)
	v := s.Get("k")
	if v < 1.9 {
		t.Fatalf("expected ~2.0, got %f", v)
	}
}
