package stale

import (
	"testing"
	"time"
)

func TestTouch_MarksPortSeen(t *testing.T) {
	tr := New(time.Minute)
	tr.Touch(8080)
	if tr.IsStale(8080) {
		t.Fatal("expected port to not be stale immediately after Touch")
	}
}

func TestSweep_MarksIdlePortStale(t *testing.T) {
	tr := New(time.Minute)
	now := time.Now()
	tr.now = func() time.Time { return now }
	tr.Touch(8080)
	// advance time beyond idle window
	tr.now = func() time.Time { return now.Add(2 * time.Minute) }
	staled := tr.Sweep()
	if len(staled) != 1 || staled[0] != 8080 {
		t.Fatalf("expected [8080] staled, got %v", staled)
	}
	if !tr.IsStale(8080) {
		t.Fatal("expected port 8080 to be stale after sweep")
	}
}

func TestSweep_DoesNotRestaleAlreadyStale(t *testing.T) {
	tr := New(time.Minute)
	now := time.Now()
	tr.now = func() time.Time { return now }
	tr.Touch(9090)
	tr.now = func() time.Time { return now.Add(2 * time.Minute) }
	first := tr.Sweep()
	second := tr.Sweep()
	if len(first) != 1 {
		t.Fatalf("expected 1 on first sweep, got %d", len(first))
	}
	if len(second) != 0 {
		t.Fatalf("expected 0 on second sweep, got %d", len(second))
	}
}

func TestTouch_ClearsStaleFlag(t *testing.T) {
	tr := New(time.Minute)
	now := time.Now()
	tr.now = func() time.Time { return now }
	tr.Touch(443)
	tr.now = func() time.Time { return now.Add(2 * time.Minute) }
	tr.Sweep()
	tr.now = func() time.Time { return now.Add(3 * time.Minute) }
	tr.Touch(443)
	if tr.IsStale(443) {
		t.Fatal("expected stale flag to be cleared after Touch")
	}
}

func TestRemove_DeletesEntry(t *testing.T) {
	tr := New(time.Minute)
	tr.Touch(22)
	tr.Remove(22)
	if tr.Len() != 0 {
		t.Fatalf("expected Len 0 after Remove, got %d", tr.Len())
	}
	if tr.IsStale(22) {
		t.Fatal("expected removed port to not be stale")
	}
}

func TestDefaultConfig(t *testing.T) {
	c := DefaultConfig()
	if c.IdleAfter <= 0 {
		t.Fatal("expected positive IdleAfter in DefaultConfig")
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("DefaultConfig.Validate: %v", err)
	}
}

func TestNewFromConfig_InvalidConfig(t *testing.T) {
	_, err := NewFromConfig(Config{IdleAfter: -1})
	if err == nil {
		t.Fatal("expected error for negative IdleAfter")
	}
}
