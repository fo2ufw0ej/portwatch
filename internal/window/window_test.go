package window

import (
	"testing"
	"time"
)

func TestAdd_IncreasesTotal(t *testing.T) {
	c := New(time.Minute, 4)
	c.Add(3)
	if got := c.Total(); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestReset_ClearsAll(t *testing.T) {
	c := New(time.Minute, 4)
	c.Add(5)
	c.Reset()
	if got := c.Total(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestTick_AdvancesBucket(t *testing.T) {
	c := New(time.Minute, 4)
	c.Add(2)
	c.Tick()
	c.Add(3)
	if got := c.Total(); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestTotal_ExcludesExpiredBuckets(t *testing.T) {
	c := New(10*time.Millisecond, 2)
	c.Add(10)
	time.Sleep(20 * time.Millisecond)
	if got := c.Total(); got != 0 {
		t.Fatalf("expected 0 after expiry, got %d", got)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Period <= 0 {
		t.Fatal("expected positive period")
	}
	if cfg.Buckets < 1 {
		t.Fatal("expected at least 1 bucket")
	}
}

func TestValidate_ZeroPeriod(t *testing.T) {
	cfg := Config{Period: 0, Buckets: 4}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero period")
	}
}

func TestValidate_ZeroBuckets(t *testing.T) {
	cfg := Config{Period: time.Minute, Buckets: 0}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero buckets")
	}
}

func TestNewFromConfig_Valid(t *testing.T) {
	cfg := DefaultConfig()
	c, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil counter")
	}
}
