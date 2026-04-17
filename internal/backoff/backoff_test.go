package backoff

import (
	"testing"
	"time"
)

func TestNext_IncreasesOverTime(t *testing.T) {
	b := New(Config{
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     10 * time.Second,
		Multiplier:      2.0,
		Jitter:          0,
	})
	prev := b.Next()
	for i := 0; i < 5; i++ {
		next := b.Next()
		if next < prev {
			t.Fatalf("expected non-decreasing intervals, got %v then %v", prev, next)
		}
		prev = next
	}
}

func TestNext_CapsAtMaxInterval(t *testing.T) {
	b := New(Config{
		InitialInterval: 1 * time.Second,
		MaxInterval:     2 * time.Second,
		Multiplier:      10.0,
		Jitter:          0,
	})
	for i := 0; i < 10; i++ {
		d := b.Next()
		if d > 2*time.Second {
			t.Fatalf("interval %v exceeded MaxInterval", d)
		}
	}
}

func TestReset_RestartsSequence(t *testing.T) {
	b := New(Config{
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     10 * time.Second,
		Multiplier:      2.0,
		Jitter:          0,
	})
	first := b.Next()
	b.Next()
	b.Next()
	b.Reset()
	if b.Attempt() != 0 {
		t.Fatalf("expected attempt 0 after reset, got %d", b.Attempt())
	}
	again := b.Next()
	if again != first {
		t.Fatalf("expected %v after reset, got %v", first, again)
	}
}

func TestValidate_InvalidMultiplier(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Multiplier = 0.5
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for Multiplier < 1")
	}
}

func TestValidate_InvalidJitter(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Jitter = 1.5
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for Jitter > 1")
	}
}

func TestNewFromConfig_Valid(t *testing.T) {
	b, err := NewFromConfig(DefaultConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil Backoff")
	}
}
