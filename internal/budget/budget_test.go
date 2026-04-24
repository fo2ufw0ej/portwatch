package budget

import (
	"testing"
	"time"
)

func defaultCfg() Config {
	return Config{
		Window:    time.Second,
		Buckets:   4,
		Threshold: 0.80,
	}
}

func TestNew_InvalidConfig(t *testing.T) {
	_, err := New(Config{Window: -1, Buckets: 4, Threshold: 0.9})
	if err == nil {
		t.Fatal("expected error for negative window")
	}
}

func TestRemaining_NoObservations_ReturnsOne(t *testing.T) {
	b, _ := New(defaultCfg())
	if got := b.Remaining(); got != 1.0 {
		t.Fatalf("expected 1.0, got %v", got)
	}
}

func TestRemaining_AllSuccesses_ReturnsOne(t *testing.T) {
	b, _ := New(defaultCfg())
	for i := 0; i < 10; i++ {
		b.Record(true)
	}
	if got := b.Remaining(); got != 1.0 {
		t.Fatalf("expected 1.0, got %v", got)
	}
}

func TestExhausted_WhenBelowThreshold(t *testing.T) {
	b, _ := New(defaultCfg()) // threshold 0.80
	// 1 success, 9 failures → 10% success rate < 80%
	b.Record(true)
	for i := 0; i < 9; i++ {
		b.Record(false)
	}
	if !b.Exhausted() {
		t.Fatal("expected budget to be exhausted")
	}
}

func TestRemaining_PartialBudget(t *testing.T) {
	b, _ := New(defaultCfg()) // threshold 0.80
	// 9 successes, 1 failure → 90% success rate
	for i := 0; i < 9; i++ {
		b.Record(true)
	}
	b.Record(false)
	got := b.Remaining()
	// remaining = (0.90 - 0.80) / (1 - 0.80) = 0.10 / 0.20 = 0.5
	if got < 0.49 || got > 0.51 {
		t.Fatalf("expected ~0.5, got %v", got)
	}
}

func TestReset_ClearsObservations(t *testing.T) {
	b, _ := New(defaultCfg())
	for i := 0; i < 10; i++ {
		b.Record(false)
	}
	b.Reset()
	if b.Exhausted() {
		t.Fatal("expected budget to be healthy after reset")
	}
}

func TestDefaultConfig_IsValid(t *testing.T) {
	if err := DefaultConfig().Validate(); err != nil {
		t.Fatalf("default config invalid: %v", err)
	}
}

func TestValidate_ThresholdOutOfRange(t *testing.T) {
	cfg := defaultCfg()
	cfg.Threshold = 1.0
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for threshold == 1.0")
	}
}
