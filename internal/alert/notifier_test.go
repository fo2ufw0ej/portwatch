package alert_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func newBufNotifier(t *testing.T) (*alert.Notifier, *bytes.Buffer) {
	t.Helper()
	var buf bytes.Buffer
	n, err := alert.NewNotifier(&buf)
	if err != nil {
		t.Fatalf("NewNotifier: %v", err)
	}
	return n, &buf
}

func TestRateLimitedNotifier_AllowsUnderBurst(t *testing.T) {
	inner, buf := newBufNotifier(t)
	rl := alert.NewRateLimitedNotifier(inner, time.Minute, 3)

	diff := scanner.Diff{Opened: []int{80}}
	for i := 0; i < 3; i++ {
		if err := rl.Notify(diff); err != nil {
			t.Fatalf("Notify error: %v", err)
		}
	}
	if !strings.Contains(buf.String(), "80") {
		t.Fatal("expected port 80 in output")
	}
}

func TestRateLimitedNotifier_SuppressesAfterBurst(t *testing.T) {
	inner, buf := newBufNotifier(t)
	rl := alert.NewRateLimitedNotifier(inner, time.Minute, 2)

	diff := scanner.Diff{Opened: []int{443}}
	rl.Notify(diff)
	rl.Notify(diff)
	buf.Reset()

	// Third call should be suppressed.
	rl.Notify(diff)
	if buf.Len() != 0 {
		t.Fatalf("expected suppressed output, got: %s", buf.String())
	}
}

func TestRateLimitedNotifier_ClosedAndOpenedIndependent(t *testing.T) {
	inner, buf := newBufNotifier(t)
	rl := alert.NewRateLimitedNotifier(inner, time.Minute, 1)

	rl.Notify(scanner.Diff{Opened: []int{22}})
	rl.Notify(scanner.Diff{Closed: []int{22}})

	out := buf.String()
	if !strings.Contains(out, "opened") && !strings.Contains(out, "OPENED") {
		t.Fatal("expected opened event in output")
	}
	if !strings.Contains(out, "closed") && !strings.Contains(out, "CLOSED") {
		t.Fatal("expected closed event in output")
	}
}

func TestRateLimitedNotifier_EmptyDiff(t *testing.T) {
	inner, buf := newBufNotifier(t)
	rl := alert.NewRateLimitedNotifier(inner, time.Minute, 1)

	if err := rl.Notify(scanner.Diff{}); err != nil {
		t.Fatalf("unexpected error on empty diff: %v", err)
	}
	if buf.Len() != 0 {
		t.Fatalf("expected no output for empty diff, got: %s", buf.String())
	}
}
