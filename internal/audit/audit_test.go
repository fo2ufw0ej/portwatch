package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/arnavsurve/portwatch/internal/audit"
)

func TestLog_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf, "text")
	err := l.Log(audit.EventPortOpened, "port opened", 9090)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "port_opened") {
		t.Errorf("expected kind in output, got: %s", out)
	}
	if !strings.Contains(out, "port=9090") {
		t.Errorf("expected port in output, got: %s", out)
	}
	if !strings.Contains(out, "port opened") {
		t.Errorf("expected message in output, got: %s", out)
	}
}

func TestLog_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf, "json")
	err := l.Log(audit.EventPortClosed, "port closed", 443)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var e audit.Event
	if err := json.Unmarshal(buf.Bytes(), &e); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if e.Kind != audit.EventPortClosed {
		t.Errorf("expected port_closed, got %s", e.Kind)
	}
	if e.Port != 443 {
		t.Errorf("expected port 443, got %d", e.Port)
	}
}

func TestLog_NoPort_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf, "text")
	_ = l.Log(audit.EventDaemonStart, "daemon started", 0)
	if strings.Contains(buf.String(), "port=") {
		t.Errorf("expected no port field when port is 0")
	}
}

func TestNew_NilWriter_DefaultsToStdout(t *testing.T) {
	l := audit.New(nil, "text")
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestLog_DefaultFormat(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf, "")
	_ = l.Log(audit.EventScanComplete, "scan done", 0)
	if buf.Len() == 0 {
		t.Error("expected output with default format")
	}
}
