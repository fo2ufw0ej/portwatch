package cadence

import (
	"testing"
	"time"
)

func defaultTracker(t *testing.T) *Tracker {
	t.Helper()
	cfg := Config{Expected: 10 * time.Second, Tolerance: 2 * time.Second, WindowSize: 5}
	tr, err := New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return tr
}

func TestObserve_FirstCall_ReturnsZero(t *testing.T) {
	tr := defaultTracker(t)
	got := tr.Observe(time.Now())
	if got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestObserve_SecondCall_ReturnsInterval(t *testing.T) {
	tr := defaultTracker(t)
	now := time.Now()
	tr.Observe(now)
	got := tr.Observe(now.Add(10 * time.Second))
	if got != 10*time.Second {
		t.Fatalf("expected 10s, got %v", got)
	}
}

func TestOnTime_WithinTolerance_ReturnsTrue(t *testing.T) {
	tr := defaultTracker(t)
	now := time.Now()
	tr.Observe(now)
	tr.Observe(now.Add(11 * time.Second)) // within ±2s of 10s
	if !tr.OnTime() {
		t.Fatal("expected OnTime to be true")
	}
}

func TestOnTime_OutsideTolerance_ReturnsFalse(t *testing.T) {
	tr := defaultTracker(t)
	now := time.Now()
	tr.Observe(now)
	tr.Observe(now.Add(30 * time.Second))
	if tr.OnTime() {
		t.Fatal("expected OnTime to be false")
	}
}

func TestOnTime_NoIntervals_ReturnsFalse(t *testing.T) {
	tr := defaultTracker(t)
	if tr.OnTime() {
		t.Fatal("expected false with no observations")
	}
}

func TestAverage_NoIntervals_ReturnsZero(t *testing.T) {
	tr := defaultTracker(t)
	if tr.Average() != 0 {
		t.Fatal("expected zero average")
	}
}

func TestAverage_MultipleIntervals(t *testing.T) {
	tr := defaultTracker(t)
	now := time.Now()
	tr.Observe(now)
	tr.Observe(now.Add(8 * time.Second))
	tr.Observe(now.Add(20 * time.Second)) // second interval = 12s
	got := tr.Average()                   // (8+12)/2 = 10s
	if got != 10*time.Second {
		t.Fatalf("expected 10s average, got %v", got)
	}
}

func TestWindowSize_LimitsHistory(t *testing.T) {
	cfg := Config{Expected: 1 * time.Second, Tolerance: 500 * time.Millisecond, WindowSize: 2}
	tr, _ := New(cfg)
	now := time.Now()
	tr.Observe(now)
	tr.Observe(now.Add(1 * time.Second))
	tr.Observe(now.Add(2 * time.Second))
	tr.Observe(now.Add(3 * time.Second))
	// window should only keep last 2 intervals
	if tr.Average() != 1*time.Second {
		t.Fatalf("unexpected average: %v", tr.Average())
	}
}

func TestReset_ClearsState(t *testing.T) {
	tr := defaultTracker(t)
	now := time.Now()
	tr.Observe(now)
	tr.Observe(now.Add(10 * time.Second))
	tr.Reset()
	if tr.Average() != 0 {
		t.Fatal("expected zero after reset")
	}
	if tr.OnTime() {
		t.Fatal("expected false after reset")
	}
}

func TestValidate_InvalidExpected(t *testing.T) {
	cfg := Config{Expected: 0, Tolerance: time.Second, WindowSize: 5}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero Expected")
	}
}

func TestValidate_NegativeTolerance(t *testing.T) {
	cfg := Config{Expected: time.Second, Tolerance: -1, WindowSize: 5}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative Tolerance")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("DefaultConfig invalid: %v", err)
	}
}
