package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Alert represents a single port change notification.
type Alert struct {
	Timestamp time.Time
	Level     Level
	Message   string
	Port      int
}

// Notifier sends alerts to a configured output.
type Notifier struct {
	out io.Writer
}

// NewNotifier creates a Notifier that writes to the given writer.
// If w is nil, os.Stdout is used.
func NewNotifier(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{out: w}
}

// Notify formats and writes alerts derived from a scanner.Diff.
// It returns all alerts generated for both newly opened and closed ports.
func (n *Notifier) Notify(diff scanner.Diff) []Alert {
	var alerts []Alert
	now := time.Now()

	for _, port := range diff.Opened {
		a := Alert{
			Timestamp: now,
			Level:     LevelAlert,
			Message:   fmt.Sprintf("Port %d newly OPENED", port),
			Port:      port,
		}
		alerts = append(alerts, a)
		n.write(a)
	}

	for _, port := range diff.Closed {
		a := Alert{
			Timestamp: now,
			Level:     LevelWarn,
			Message:   fmt.Sprintf("Port %d unexpectedly CLOSED", port),
			Port:      port,
		}
		alerts = append(alerts, a)
		n.write(a)
	}

	return alerts
}

// NotifyInfo writes an informational alert with the given message.
// This is useful for surfacing non-critical status updates (e.g. scan start/stop).
func (n *Notifier) NotifyInfo(message string) Alert {
	a := Alert{
		Timestamp: time.Now(),
		Level:     LevelInfo,
		Message:   message,
		Port:      0,
	}
	n.write(a)
	return a
}

func (n *Notifier) write(a Alert) {
	fmt.Fprintf(n.out, "[%s] %s %s\n",
		a.Timestamp.Format(time.RFC3339),
		a.Level,
		a.Message,
	)
}
