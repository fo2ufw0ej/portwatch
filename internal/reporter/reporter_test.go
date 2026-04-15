package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/reporter"
	"github.com/user/portwatch/internal/scanner"
)

func TestWriteSummary_Text(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	if err := r.WriteSummary([]int{80, 443, 8080}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Open ports (3)") {
		t.Errorf("expected port count in output, got: %s", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port 8080 in output, got: %s", out)
	}
}

func TestWriteSummary_JSON(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatJSON)
	if err := r.WriteSummary([]int{22, 80}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"open_port_count":2`) {
		t.Errorf("expected count field in JSON output, got: %s", out)
	}
	if !strings.Contains(out, "22") {
		t.Errorf("expected port 22 in JSON output, got: %s", out)
	}
}

func TestWriteSummary_NilWriter_DefaultsToStdout(t *testing.T) {
	// Should not panic when out is nil (defaults to os.Stdout).
	r := reporter.New(nil, reporter.FormatText)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}

func TestWriteDiff_WithChanges(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	diff := scanner.Diff{Opened: []int{9090}, Closed: []int{8080}}
	if err := r.WriteDiff(diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "9090") {
		t.Errorf("expected opened port in output, got: %s", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected closed port in output, got: %s", out)
	}
}

func TestWriteDiff_NoChanges_WritesNothing(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	diff := scanner.Diff{}
	if err := r.WriteDiff(diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diff, got: %s", buf.String())
	}
}
