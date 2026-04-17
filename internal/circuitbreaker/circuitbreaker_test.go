package circuitbreaker

import (
	"testing"
	"time"
)

func TestAllow_ClosedByDefault(t *testing.T) {
	b := New(3, 50*time.Millisecond)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestOpensAfterMaxFailures(t *testing.T) {
	b := New(3, 50*time.Millisecond)
	b.RecordFailure()
	b.RecordFailure()
	if b.State() != StateClosed {
		t.Fatal("should still be closed")
	}
	b.RecordFailure()
	if b.State() != StateOpen {
		t.Fatal("should be open after max failures")
	}
	if err := b.Allow(); err != ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestHalfOpenAfterReset(t *testing.T) {
	b := New(1, 20*time.Millisecond)
	b.RecordFailure()
	time.Sleep(30 * time.Millisecond)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil in half-open, got %v", err)
	}
	if b.State() != StateHalfOpen {
		t.Fatal("expected half-open state")
	}
}

func TestRecordSuccess_ResetsClosed(t *testing.T) {
	b := New(2, 50*time.Millisecond)
	b.RecordFailure()
	b.RecordFailure()
	b.RecordSuccess()
	if b.State() != StateClosed {
		t.Fatal("expected closed after success")
	}
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestDefaultConfig(t *testing.T) {
	c := DefaultConfig()
	if err := c.Validate(); err != nil {
		t.Fatalf("default config invalid: %v", err)
	}
}

func TestNewFromConfig_Invalid(t *testing.T) {
	c := Config{MaxFailures: 0, ResetAfter: time.Second}
	if _, err := NewFromConfig(c); err == nil {
		t.Fatal("expected error for zero max_failures")
	}
}

func TestNewFromConfig_Valid(t *testing.T) {
	c := DefaultConfig()
	b, err := NewFromConfig(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil breaker")
	}
}
