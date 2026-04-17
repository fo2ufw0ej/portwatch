package sampler

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Interval <= 0 {
		t.Fatal("expected positive interval")
	}
	if cfg.Jitter < 0 {
		t.Fatal("expected non-negative jitter")
	}
}

func TestValidate_InvalidInterval(t *testing.T) {
	cfg := Config{Interval: 0, Jitter: 0}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestValidate_NegativeJitter(t *testing.T) {
	cfg := Config{Interval: time.Second, Jitter: -1}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative jitter")
	}
}

func TestNew_InvalidConfig(t *testing.T) {
	_, err := New(Config{Interval: -1})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_FiresCallback(t *testing.T) {
	cfg := Config{Interval: 20 * time.Millisecond, Jitter: 0}
	s, err := New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var count int64
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		s.Run(ctx, func() { atomic.AddInt64(&count, 1) })
		close(done)
	}()

	<-done
	got := atomic.LoadInt64(&count)
	if got < 2 {
		t.Fatalf("expected at least 2 callbacks, got %d", got)
	}
}

func TestRun_StopsOnCancel(t *testing.T) {
	cfg := Config{Interval: 10 * time.Millisecond, Jitter: 0}
	s, err := New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	done := make(chan struct{})
	go func() {
		s.Run(ctx, func() {})
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Run did not stop after context cancellation")
	}
}
