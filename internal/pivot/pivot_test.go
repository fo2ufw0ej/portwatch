package pivot

import (
	"testing"
	"time"
)

func TestRecord_AndCount(t *testing.T) {
	tr := New(10 * time.Second)
	tr.Record(80)
	tr.Record(80)
	tr.Record(443)

	if got := tr.Count(80); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
	if got := tr.Count(443); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestCount_MissingPort_ReturnsZero(t *testing.T) {
	tr := New(10 * time.Second)
	if got := tr.Count(9999); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestPrune_RemovesExpiredEvents(t *testing.T) {
	tr := New(5 * time.Second)

	base := time.Now()
	tr.now = func() time.Time { return base }
	tr.Record(22)
	tr.Record(22)

	// advance past window
	tr.now = func() time.Time { return base.Add(10 * time.Second) }

	if got := tr.Count(22); got != 0 {
		t.Fatalf("expected 0 after expiry, got %d", got)
	}
}

func TestTop_ReturnsSortedByCount(t *testing.T) {
	tr := New(30 * time.Second)
	for i := 0; i < 3; i++ {
		tr.Record(80)
	}
	for i := 0; i < 5; i++ {
		tr.Record(443)
	}
	tr.Record(22)

	top := tr.Top(2)
	if len(top) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(top))
	}
	if top[0].Port != 443 || top[0].Transitions != 5 {
		t.Errorf("expected 443 first with 5 transitions, got port=%d count=%d", top[0].Port, top[0].Transitions)
	}
	if top[1].Port != 80 || top[1].Transitions != 3 {
		t.Errorf("expected 80 second with 3 transitions, got port=%d count=%d", top[1].Port, top[1].Transitions)
	}
}

func TestTop_NLimit_ReturnsAll(t *testing.T) {
	tr := New(30 * time.Second)
	tr.Record(8080)
	tr.Record(3000)

	top := tr.Top(0)
	if len(top) != 2 {
		t.Fatalf("expected 2 entries with n=0, got %d", len(top))
	}
}

func TestReset_ClearsAll(t *testing.T) {
	tr := New(30 * time.Second)
	tr.Record(80)
	tr.Record(443)
	tr.Reset()

	if got := tr.Count(80); got != 0 {
		t.Errorf("expected 0 after reset, got %d", got)
	}
	if top := tr.Top(10); len(top) != 0 {
		t.Errorf("expected empty top after reset, got %d entries", len(top))
	}
}
