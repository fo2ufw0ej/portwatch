package watermark_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/watermark"
)

func TestNew_NoObservations(t *testing.T) {
	w := watermark.New()
	m := w.Snapshot()
	if m.Observed != 0 {
		t.Fatalf("expected 0 observations, got %d", m.Observed)
	}
}

func TestObserve_SetsInitialHighAndLow(t *testing.T) {
	w := watermark.New()
	w.Observe(5)
	m := w.Snapshot()
	if m.High != 5 {
		t.Errorf("expected High=5, got %d", m.High)
	}
	if m.Low != 5 {
		t.Errorf("expected Low=5, got %d", m.Low)
	}
	if m.Observed != 1 {
		t.Errorf("expected Observed=1, got %d", m.Observed)
	}
}

func TestObserve_UpdatesHigh(t *testing.T) {
	w := watermark.New()
	w.Observe(3)
	w.Observe(10)
	w.Observe(7)
	m := w.Snapshot()
	if m.High != 10 {
		t.Errorf("expected High=10, got %d", m.High)
	}
	if m.Low != 3 {
		t.Errorf("expected Low=3, got %d", m.Low)
	}
	if m.Observed != 3 {
		t.Errorf("expected Observed=3, got %d", m.Observed)
	}
}

func TestObserve_UpdatesLow(t *testing.T) {
	w := watermark.New()
	w.Observe(8)
	w.Observe(2)
	m := w.Snapshot()
	if m.Low != 2 {
		t.Errorf("expected Low=2, got %d", m.Low)
	}
}

func TestReset_ClearsMarks(t *testing.T) {
	w := watermark.New()
	w.Observe(5)
	w.Observe(12)
	w.Reset()
	m := w.Snapshot()
	if m.Observed != 0 {
		t.Errorf("expected 0 observations after reset, got %d", m.Observed)
	}
	if m.High != 0 || m.Low != 0 {
		t.Errorf("expected zero marks after reset")
	}
}

func TestWrite_NoObservations(t *testing.T) {
	w := watermark.New()
	var buf bytes.Buffer
	w.Write(&buf)
	if !strings.Contains(buf.String(), "no observations") {
		t.Errorf("expected 'no observations' message, got: %s", buf.String())
	}
}

func TestWrite_WithObservations(t *testing.T) {
	w := watermark.New()
	w.Observe(4)
	w.Observe(9)
	var buf bytes.Buffer
	w.Write(&buf)
	out := buf.String()
	if !strings.Contains(out, "high=9") {
		t.Errorf("expected high=9 in output, got: %s", out)
	}
	if !strings.Contains(out, "low=4") {
		t.Errorf("expected low=4 in output, got: %s", out)
	}
}

func TestWrite_NilWriter_DefaultsToStdout(t *testing.T) {
	w := watermark.New()
	w.Observe(1)
	// Should not panic
	w.Write(nil)
}
