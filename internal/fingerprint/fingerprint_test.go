package fingerprint_test

import (
	"testing"

	"github.com/user/portwatch/internal/fingerprint"
)

func TestCompute_OrderIndependent(t *testing.T) {
	a := fingerprint.Compute([]int{80, 443, 22})
	b := fingerprint.Compute([]int{443, 22, 80})
	if a != b {
		t.Fatalf("expected same fingerprint for same ports in different order, got %s vs %s", a, b)
	}
}

func TestCompute_EmptyPorts(t *testing.T) {
	f := fingerprint.Compute([]int{})
	if f != fingerprint.Empty {
		t.Fatalf("expected Empty fingerprint, got %s", f)
	}
}

func TestCompute_NilPorts(t *testing.T) {
	f := fingerprint.Compute(nil)
	if f != fingerprint.Empty {
		t.Fatalf("expected Empty fingerprint for nil slice, got %s", f)
	}
}

func TestCompute_DifferentPorts(t *testing.T) {
	a := fingerprint.Compute([]int{80, 443})
	b := fingerprint.Compute([]int{80, 8080})
	if fingerprint.Equal(a, b) {
		t.Fatal("expected different fingerprints for different port sets")
	}
}

func TestChanged_DetectsChange(t *testing.T) {
	prev := fingerprint.Compute([]int{22, 80})
	if !fingerprint.Changed(prev, []int{22, 80, 443}) {
		t.Fatal("expected Changed to return true when port added")
	}
}

func TestChanged_NoChange(t *testing.T) {
	prev := fingerprint.Compute([]int{22, 80})
	if fingerprint.Changed(prev, []int{80, 22}) {
		t.Fatal("expected Changed to return false for same ports")
	}
}

func TestShort_TruncatesTo12(t *testing.T) {
	f := fingerprint.Compute([]int{80})
	short := f.Short()
	if len(short) != 12 {
		t.Fatalf("expected Short() length 12, got %d", len(short))
	}
}

func TestString_FullHex(t *testing.T) {
	f := fingerprint.Compute([]int{443})
	if len(f.String()) != 64 {
		t.Fatalf("expected full 64-char hex string, got length %d", len(f.String()))
	}
}
