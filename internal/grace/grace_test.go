package grace_test

import (
	"testing"
	"time"

	"github.com/densestvoid/portwatch/internal/grace"
)

func TestShutdown_AllFinish(t *testing.T) {
	cfg := grace.DefaultConfig()
	coord, err := grace.New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	coord.Register()
	go func() {
		defer coord.Done()
		<-coord.Context().Done()
	}()

	ok := coord.Shutdown()
	if !ok {
		t.Error("expected clean shutdown")
	}
}

func TestShutdown_Timeout(t *testing.T) {
	cfg := grace.Config{Timeout: 50 * time.Millisecond}
	coord, err := grace.New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	coord.Register()
	// intentionally never call Done

	ok := coord.Shutdown()
	if ok {
		t.Error("expected timeout, got clean shutdown")
	}
}

func TestContext_CancelledOnShutdown(t *testing.T) {
	cfg := grace.DefaultConfig()
	coord, _ := grace.New(cfg)

	select {
	case <-coord.Context().Done():
		t.Fatal("context should not be cancelled yet")
	default:
	}

	go coord.Shutdown()

	select {
	case <-coord.Context().Done():
		// expected
	case <-time.After(time.Second):
		t.Fatal("context was not cancelled after Shutdown")
	}
}

func TestNew_InvalidConfig(t *testing.T) {
	_, err := grace.New(grace.Config{Timeout: 0})
	if err == nil {
		t.Error("expected error for zero timeout")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := grace.DefaultConfig()
	if cfg.Timeout <= 0 {
		t.Error("expected positive default timeout")
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("default config invalid: %v", err)
	}
}
