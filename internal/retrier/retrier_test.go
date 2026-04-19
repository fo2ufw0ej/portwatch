package retrier

import (
	"context"
	"errors"
	"testing"
	"time"
)

func fastConfig(attempts int) Config {
	return Config{
		MaxAttempts:     attempts,
		InitialInterval: time.Millisecond,
		MaxInterval:     5 * time.Millisecond,
		Multiplier:      2.0,
	}
}

func TestDo_SucceedsFirstAttempt(t *testing.T) {
	r, _ := New(fastConfig(3))
	calls := 0
	err := r.Do(context.Background(), func() error {
		calls++
		return nil
	})
	if err != nil || calls != 1 {
		t.Fatalf("expected 1 call and no error, got calls=%d err=%v", calls, err)
	}
}

func TestDo_RetriesOnError(t *testing.T) {
	r, _ := New(fastConfig(3))
	calls := 0
	sentinel := errors.New("transient")
	err := r.Do(context.Background(), func() error {
		calls++
		if calls < 3 {
			return sentinel
		}
		return nil
	})
	if err != nil || calls != 3 {
		t.Fatalf("expected 3 calls and no error, got calls=%d err=%v", calls, err)
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	r, _ := New(fastConfig(2))
	sentinel := errors.New("always fails")
	calls := 0
	err := r.Do(context.Background(), func() error {
		calls++
		return sentinel
	})
	if !errors.Is(err, sentinel) || calls != 2 {
		t.Fatalf("expected sentinel after 2 calls, got calls=%d err=%v", calls, err)
	}
}

func TestDo_PermanentError_NoRetry(t *testing.T) {
	r, _ := New(fastConfig(5))
	calls := 0
	err := r.Do(context.Background(), func() error {
		calls++
		return Permanent(errors.New("fatal"))
	})
	if err == nil || calls != 1 {
		t.Fatalf("expected immediate stop, got calls=%d err=%v", calls, err)
	}
}

func TestDo_ContextCancelled(t *testing.T) {
	r, _ := New(fastConfig(10))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := r.Do(ctx, func() error { return errors.New("x") })
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestValidate_InvalidMultiplier(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Multiplier = 0.5
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for multiplier < 1")
	}
}

func TestDefaultConfig_IsValid(t *testing.T) {
	if err := DefaultConfig().Validate(); err != nil {
		t.Fatalf("default config invalid: %v", err)
	}
}
