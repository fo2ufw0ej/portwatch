package daemon_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/config"
	"github.com/yourorg/portwatch/internal/daemon"
	"github.com/yourorg/portwatch/internal/state"
)

func buildConfig(t *testing.T) *config.Config {
	t.Helper()
	cfg := config.DefaultConfig()
	// Restrict to a high unprivileged range unlikely to have open ports in CI.
	cfg.PortRange.Start = 59000
	cfg.PortRange.End = 59010
	cfg.Interval = 50 * time.Millisecond
	cfg.AlertOutput = "stdout"
	return cfg
}

func TestDaemon_RunAndStop(t *testing.T) {
	cfg := buildConfig(t)
	storePath := filepath.Join(t.TempDir(), "state.json")

	d, err := daemon.New(cfg, storePath)
	if err != nil {
		t.Fatalf("daemon.New: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	err = d.Run(ctx)
	// context.DeadlineExceeded is the expected exit reason.
	if err != context.DeadlineExceeded {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}

func TestDaemon_PersistsState(t *testing.T) {
	cfg := buildConfig(t)
	storePath := filepath.Join(t.TempDir(), "state.json")

	d, err := daemon.New(cfg, storePath)
	if err != nil {
		t.Fatalf("daemon.New: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()
	_ = d.Run(ctx)

	store := state.NewStore(storePath)
	snap, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if snap.RecordedAt.IsZero() {
		t.Error("expected non-zero RecordedAt after daemon run")
	}
}
