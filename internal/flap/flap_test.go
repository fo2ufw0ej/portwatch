package flap_test

import (
	"testing"
	"time"

	"github.com/wiliamsouza/portwatch/internal/flap"
)

func TestDefaultConfig(t *testing.T) {
	cfg := flap.DefaultConfig()
	if cfg.Window <= 0 {
		t.Fatalf("expected positive Window, got %v", cfg.Window)
	}
	if cfg.Threshold < 2 {
		t.Fatalf("expected Threshold >= 2, got %d", cfg.Threshold)
	}
}

func TestValidate_InvalidWindow(t *testing.T) {
	cfg := flap.DefaultConfig()
	cfg.Window = 0
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero Window")
	}
}

func TestValidate_InvalidThreshold(t *testing.T) {
	cfg := flap.DefaultConfig()
	cfg.Threshold = 1
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for Threshold < 2")
	}
}

func TestNew_InvalidConfig(t *testing.T) {
	cfg := flap.DefaultConfig()
	cfg.Window = -1
	_, err := flap.New(cfg)
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}

func TestFlapping_BelowThreshold(t *testing.T) {
	cfg := flap.DefaultConfig()
	cfg.Window = time.Minute
	cfg.Threshold = 4
	d, _ := flap.New(cfg)

	now := time.Now()
	d.Record(8080, now)
	d.Record(8080, now.Add(5*time.Second))
	d.Record(8080, now.Add(10*time.Second))

	if d.Flapping(8080, now.Add(15*time.Second)) {
		t.Fatal("expected not flapping below threshold")
	}
}

func TestFlapping_AtThreshold(t *testing.T) {
	cfg := flap.DefaultConfig()
	cfg.Window = time.Minute
	cfg.Threshold = 4
	d, _ := flap.New(cfg)

	now := time.Now()
	for i := 0; i < 4; i++ {
		d.Record(9090, now.Add(time.Duration(i)*5*time.Second))
	}

	if !d.Flapping(9090, now.Add(30*time.Second)) {
		t.Fatal("expected flapping at threshold")
	}
}

func TestFlapping_EventsExpireOutsideWindow(t *testing.T) {
	cfg := flap.DefaultConfig()
	cfg.Window = 30 * time.Second
	cfg.Threshold = 3
	d, _ := flap.New(cfg)

	base := time.Now()
	// Record 3 events far in the past — they should be pruned.
	for i := 0; i < 3; i++ {
		d.Record(443, base.Add(time.Duration(i)*time.Second))
	}

	now := base.Add(2 * time.Minute)
	if d.Flapping(443, now) {
		t.Fatal("expected old events to be pruned")
	}
}

func TestReset_ClearsEvents(t *testing.T) {
	cfg := flap.DefaultConfig()
	cfg.Window = time.Minute
	cfg.Threshold = 2
	d, _ := flap.New(cfg)

	now := time.Now()
	d.Record(22, now)
	d.Record(22, now.Add(time.Second))
	d.Reset(22)

	if d.Flapping(22, now.Add(2*time.Second)) {
		t.Fatal("expected no flapping after reset")
	}
}

func TestFlapping_IndependentPorts(t *testing.T) {
	cfg := flap.DefaultConfig()
	cfg.Window = time.Minute
	cfg.Threshold = 2
	d, _ := flap.New(cfg)

	now := time.Now()
	d.Record(80, now)
	d.Record(80, now.Add(time.Second))

	if d.Flapping(443, now.Add(2*time.Second)) {
		t.Fatal("port 443 should not be affected by port 80 events")
	}
}
