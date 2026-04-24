package burst

import (
	"testing"
	"time"
)

func newTracker(t *testing.T, window time.Duration, ceiling int) *Tracker {
	t.Helper()
	tr, err := New(Config{Window: window, Ceiling: ceiling})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return tr
}

func TestRecord_IncreasesCount(t *testing.T) {
	tr := newTracker(t, time.Minute, 10)
	tr.Record("k")
	tr.Record("k")
	if got := tr.Count("k"); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestCount_EmptyKey_ReturnsZero(t *testing.T) {
	tr := newTracker(t, time.Minute, 5)
	if got := tr.Count("missing"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestSpiking_BelowCeiling_ReturnsFalse(t *testing.T) {
	tr := newTracker(t, time.Minute, 5)
	for i := 0; i < 5; i++ {
		tr.Record("k")
	}
	if tr.Spiking("k") {
		t.Fatal("expected not spiking at ceiling")
	}
}

func TestSpiking_AboveCeiling_ReturnsTrue(t *testing.T) {
	tr := newTracker(t, time.Minute, 3)
	for i := 0; i < 4; i++ {
		tr.Record("k")
	}
	if !tr.Spiking("k") {
		t.Fatal("expected spiking above ceiling")
	}
}

func TestReset_ClearsCount(t *testing.T) {
	tr := newTracker(t, time.Minute, 5)
	tr.Record("k")
	tr.Record("k")
	tr.Reset("k")
	if got := tr.Count("k"); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestCount_PrunesExpiredEvents(t *testing.T) {
	tr := newTracker(t, 50*time.Millisecond, 10)
	tr.Record("k")
	tr.Record("k")
	time.Sleep(80 * time.Millisecond)
	if got := tr.Count("k"); got != 0 {
		t.Fatalf("expected 0 after window expiry, got %d", got)
	}
}

func TestIndependentKeys(t *testing.T) {
	tr := newTracker(t, time.Minute, 5)
	tr.Record("a")
	tr.Record("b")
	tr.Record("b")
	if tr.Count("a") != 1 {
		t.Fatalf("expected 1 for key a")
	}
	if tr.Count("b") != 2 {
		t.Fatalf("expected 2 for key b")
	}
}
