package quorum

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Threshold != 2 {
		t.Fatalf("expected threshold 2, got %d", cfg.Threshold)
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestValidate_ZeroThreshold(t *testing.T) {
	cfg := Config{Threshold: 0}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero threshold")
	}
}

func TestNew_InvalidConfig(t *testing.T) {
	_, err := New(Config{Threshold: -1})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestConfirm_BelowThreshold_NotForwarded(t *testing.T) {
	q, _ := New(Config{Threshold: 3})
	diff := scanner.Diff{Opened: []int{8080}}

	// Two observations — threshold is 3, should not forward yet.
	for i := 0; i < 2; i++ {
		out := q.Confirm(diff)
		if len(out.Opened) != 0 {
			t.Fatalf("iter %d: expected no output, got %v", i, out.Opened)
		}
	}
}

func TestConfirm_AtThreshold_Forwarded(t *testing.T) {
	q, _ := New(Config{Threshold: 2})
	diff := scanner.Diff{Opened: []int{9000}}

	if out := q.Confirm(diff); len(out.Opened) != 0 {
		t.Fatalf("first call should not forward, got %v", out.Opened)
	}
	out := q.Confirm(diff)
	if len(out.Opened) != 1 || out.Opened[0] != 9000 {
		t.Fatalf("expected port 9000 forwarded, got %v", out.Opened)
	}
}

func TestConfirm_CounterResetsOnAbsence(t *testing.T) {
	q, _ := New(Config{Threshold: 3})
	diff := scanner.Diff{Opened: []int{443}}

	q.Confirm(diff) // count = 1
	q.Confirm(scanner.Diff{}) // absent → reset
	q.Confirm(diff) // count = 1 again
	out := q.Confirm(diff) // count = 2, still below 3
	if len(out.Opened) != 0 {
		t.Fatalf("expected no output after reset, got %v", out.Opened)
	}
}

func TestConfirm_ClosedPort(t *testing.T) {
	q, _ := New(Config{Threshold: 2})
	diff := scanner.Diff{Closed: []int{22}}

	q.Confirm(diff)
	out := q.Confirm(diff)
	if len(out.Closed) != 1 || out.Closed[0] != 22 {
		t.Fatalf("expected port 22 closed forwarded, got %v", out.Closed)
	}
}

func TestReset_ClearsCounters(t *testing.T) {
	q, _ := New(Config{Threshold: 2})
	diff := scanner.Diff{Opened: []int{80}}
	q.Confirm(diff)
	q.Reset()
	// After reset the counter restarts; second call should not forward.
	q.Confirm(diff)
	out := q.Confirm(diff) // now at threshold again
	if len(out.Opened) != 1 {
		t.Fatalf("expected port forwarded after reset cycle, got %v", out.Opened)
	}
}
