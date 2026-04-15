package metrics_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

func makeSnapshot() metrics.Snapshot {
	return metrics.Snapshot{
		ScansTotal:    7,
		AlertsTotal:   2,
		LastScanAt:    time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		LastScanPorts: 15,
		Uptime:        90 * time.Second,
	}
}

func TestWrite_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	s := makeSnapshot()
	if err := metrics.Write(&buf, s, metrics.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"Scans total", "7", "Alerts total", "2", "1m30s"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output:\n%s", want, out)
		}
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	s := makeSnapshot()
	if err := metrics.Write(&buf, s, metrics.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got metrics.Snapshot
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got.ScansTotal != s.ScansTotal {
		t.Errorf("scans_total mismatch: got %d want %d", got.ScansTotal, s.ScansTotal)
	}
	if got.AlertsTotal != s.AlertsTotal {
		t.Errorf("alerts_total mismatch: got %d want %d", got.AlertsTotal, s.AlertsTotal)
	}
}

func TestWrite_NilWriter_DefaultsToStdout(t *testing.T) {
	// Should not panic when writer is nil.
	s := makeSnapshot()
	if err := metrics.Write(nil, s, metrics.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
