package snapshot_test

import (
	"testing"
	"time"

	"github.com/iamkevindemarchi/portwatch/internal/snapshot"
)

func TestNew_SortsPorts(t *testing.T) {
	s := snapshot.New("localhost", []int{443, 80, 22, 8080})
	expected := []int{22, 80, 443, 8080}
	for i, p := range s.Ports {
		if p != expected[i] {
			t.Errorf("expected port %d at index %d, got %d", expected[i], i, p)
		}
	}
}

func TestNew_SetsTimestamp(t *testing.T) {
	before := time.Now()
	s := snapshot.New("localhost", []int{80})
	after := time.Now()
	if s.Timestamp.Before(before) || s.Timestamp.After(after) {
		t.Errorf("timestamp %v not within expected range [%v, %v]", s.Timestamp, before, after)
	}
}

func TestEqual_SamePorts(t *testing.T) {
	a := snapshot.New("localhost", []int{80, 443})
	b := snapshot.New("localhost", []int{443, 80})
	if !a.Equal(b) {
		t.Error("expected snapshots with same ports to be equal")
	}
}

func TestEqual_DifferentPorts(t *testing.T) {
	a := snapshot.New("localhost", []int{80, 443})
	b := snapshot.New("localhost", []int{80, 8080})
	if a.Equal(b) {
		t.Error("expected snapshots with different ports to not be equal")
	}
}

func TestEqual_DifferentLengths(t *testing.T) {
	a := snapshot.New("localhost", []int{80})
	b := snapshot.New("localhost", []int{80, 443})
	if a.Equal(b) {
		t.Error("expected snapshots with different lengths to not be equal")
	}
}

func TestContains_Found(t *testing.T) {
	s := snapshot.New("localhost", []int{22, 80, 443})
	if !s.Contains(80) {
		t.Error("expected snapshot to contain port 80")
	}
}

func TestContains_NotFound(t *testing.T) {
	s := snapshot.New("localhost", []int{22, 80, 443})
	if s.Contains(8080) {
		t.Error("expected snapshot to not contain port 8080")
	}
}

func TestSummary_WithPorts(t *testing.T) {
	s := snapshot.New("myhost", []int{80, 443})
	sum := s.Summary()
	if len(sum) == 0 {
		t.Error("expected non-empty summary")
	}
	for _, want := range []string{"myhost", "80", "443", "2 open port(s)"} {
		if !containsStr(sum, want) {
			t.Errorf("summary %q missing expected substring %q", sum, want)
		}
	}
}

func TestSummary_NoPorts(t *testing.T) {
	s := snapshot.New("myhost", []int{})
	sum := s.Summary()
	if !containsStr(sum, "no open ports") {
		t.Errorf("expected 'no open ports' in summary, got: %q", sum)
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
