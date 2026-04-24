package trend_test

import (
	"testing"
	"time"

	"github.com/montokapro/portwatch/internal/trend"
)

func TestDirection_Stable_WhenNoSamples(t *testing.T) {
	tr := trend.New(time.Minute)
	if got := tr.Direction(); got != trend.Stable {
		t.Fatalf("expected Stable, got %s", got)
	}
}

func TestDirection_Rising_WhenPositiveSum(t *testing.T) {
	tr := trend.New(time.Minute)
	tr.Record(3)
	tr.Record(1)
	if got := tr.Direction(); got != trend.Rising {
		t.Fatalf("expected Rising, got %s", got)
	}
}

func TestDirection_Falling_WhenNegativeSum(t *testing.T) {
	tr := trend.New(time.Minute)
	tr.Record(-2)
	if got := tr.Direction(); got != trend.Falling {
		t.Fatalf("expected Falling, got %s", got)
	}
}

func TestDirection_Stable_WhenSumIsZero(t *testing.T) {
	tr := trend.New(time.Minute)
	tr.Record(2)
	tr.Record(-2)
	if got := tr.Direction(); got != trend.Stable {
		t.Fatalf("expected Stable, got %s", got)
	}
}

func TestSamples_ReturnsCopy(t *testing.T) {
	tr := trend.New(time.Minute)
	tr.Record(1)
	tr.Record(-1)
	samples := tr.Samples()
	if len(samples) != 2 {
		t.Fatalf("expected 2 samples, got %d", len(samples))
	}
}

func TestReset_ClearsSamples(t *testing.T) {
	tr := trend.New(time.Minute)
	tr.Record(5)
	tr.Reset()
	if got := tr.Direction(); got != trend.Stable {
		t.Fatalf("expected Stable after reset, got %s", got)
	}
	if n := len(tr.Samples()); n != 0 {
		t.Fatalf("expected 0 samples after reset, got %d", n)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := trend.DefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config should be valid: %v", err)
	}
}

func TestValidate_ZeroWindow(t *testing.T) {
	cfg := trend.Config{Window: 0}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewFromConfig_Valid(t *testing.T) {
	cfg := trend.DefaultConfig()
	tr, err := trend.NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil Tracker")
	}
}

func TestNewFromConfig_Invalid(t *testing.T) {
	cfg := trend.Config{Window: -1 * time.Second}
	_, err := trend.NewFromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}
