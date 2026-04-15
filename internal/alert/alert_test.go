package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func TestNotify_OpenedPorts(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewNotifier(&buf)

	diff := scanner.Diff{
		Opened: []int{8080, 9090},
		Closed: []int{},
	}

	alerts := n.Notify(diff)

	if len(alerts) != 2 {
		t.Fatalf("expected 2 alerts, got %d", len(alerts))
	}

	for _, a := range alerts {
		if a.Level != alert.LevelAlert {
			t.Errorf("expected level ALERT for opened port, got %s", a.Level)
		}
	}

	out := buf.String()
	if !strings.Contains(out, "8080") || !strings.Contains(out, "9090") {
		t.Errorf("output missing port numbers: %s", out)
	}
}

func TestNotify_ClosedPorts(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewNotifier(&buf)

	diff := scanner.Diff{
		Opened: []int{},
		Closed: []int{3000},
	}

	alerts := n.Notify(diff)

	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Level != alert.LevelWarn {
		t.Errorf("expected level WARN for closed port, got %s", alerts[0].Level)
	}
	if alerts[0].Port != 3000 {
		t.Errorf("expected port 3000, got %d", alerts[0].Port)
	}
}

func TestNotify_NoDiff(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewNotifier(&buf)

	diff := scanner.Diff{}
	alerts := n.Notify(diff)

	if len(alerts) != 0 {
		t.Errorf("expected no alerts for empty diff, got %d", len(alerts))
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diff")
	}
}

func TestNewNotifier_DefaultsToStdout(t *testing.T) {
	n := alert.NewNotifier(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
