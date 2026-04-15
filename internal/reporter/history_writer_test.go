package reporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/history"
	"github.com/yourorg/portwatch/internal/reporter"
)

func sampleEntries() []history.Entry {
	return []history.Entry{
		{
			Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Opened:    []int{80, 443},
			Closed:    []int{22},
		},
		{
			Timestamp: time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC),
			Opened:    nil,
			Closed:    []int{8080},
		},
	}
}

func TestWriteHistory_Text(t *testing.T) {
	var buf bytes.Buffer
	if err := reporter.WriteHistory(&buf, sampleEntries(), "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "TIMESTAMP") {
		t.Error("expected header row")
	}
	if !strings.Contains(out, "80") {
		t.Error("expected port 80 in output")
	}
	if !strings.Contains(out, "22") {
		t.Error("expected port 22 in output")
	}
}

func TestWriteHistory_JSON(t *testing.T) {
	var buf bytes.Buffer
	if err := reporter.WriteHistory(&buf, sampleEntries(), "json"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var entries []history.Entry
	if err := json.Unmarshal(buf.Bytes(), &entries); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestWriteHistory_Empty(t *testing.T) {
	var buf bytes.Buffer
	if err := reporter.WriteHistory(&buf, nil, "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No history") {
		t.Error("expected empty-history message")
	}
}

func TestWriteHistory_NilWriter_DefaultsToStdout(t *testing.T) {
	// Should not panic when writer is nil
	if err := reporter.WriteHistory(nil, sampleEntries(), "text"); err != nil {
		t.Fatalf("unexpected error with nil writer: %v", err)
	}
}
