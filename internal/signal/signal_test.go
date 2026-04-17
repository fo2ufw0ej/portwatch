package signal_test

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/user/portwatch/internal/signal"
)

func TestHandler_CancelsOnSIGINT(t *testing.T) {
	h, ctx := signal.New(context.Background())
	defer h.Stop()

	// Send SIGINT to the current process.
	syscall.Kill(syscall.Getpid(), syscall.SIGINT) //nolint:errcheck

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(2 * time.Second):
		t.Fatal("context was not cancelled after SIGINT")
	}
}

func TestHandler_StopCancelsContext(t *testing.T) {
	h, ctx := signal.New(context.Background())
	h.Stop()

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(time.Second):
		t.Fatal("context was not cancelled after Stop")
	}
}

func TestHandler_StopIsIdempotent(t *testing.T) {
	h, _ := signal.New(context.Background())
	// Calling Stop multiple times should not panic.
	h.Stop()
	h.Stop()
}

func TestDefaultConfig(t *testing.T) {
	cfg := signal.DefaultConfig()
	if len(cfg.Signals) == 0 {
		t.Fatal("expected default signals to be non-empty")
	}
}

func TestValidate_NoSignals(t *testing.T) {
	cfg := signal.Config{Signals: []string{}}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty signals")
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := signal.DefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
