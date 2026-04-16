// Package audit provides structured audit logging for portwatch events.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// EventKind describes the type of audit event.
type EventKind string

const (
	EventScanComplete EventKind = "scan_complete"
	EventPortOpened   EventKind = "port_opened"
	EventPortClosed   EventKind = "port_closed"
	EventDaemonStart  EventKind = "daemon_start"
	EventDaemonStop   EventKind = "daemon_stop"
)

// Event is a single audit log entry.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Kind      EventKind `json:"kind"`
	Message   string    `json:"message"`
	Port      int       `json:"port,omitempty"`
}

// Logger writes audit events to a destination.
type Logger struct {
	w      io.Writer
	format string
}

// New creates a new audit Logger. format is "text" or "json".
func New(w io.Writer, format string) *Logger {
	if w == nil {
		w = os.Stdout
	}
	if format == "" {
		format = "text"
	}
	return &Logger{w: w, format: format}
}

// Log writes an audit event.
func (l *Logger) Log(kind EventKind, message string, port int) error {
	e := Event{
		Timestamp: time.Now().UTC(),
		Kind:      kind,
		Message:   message,
		Port:      port,
	}
	if l.format == "json" {
		return json.NewEncoder(l.w).Encode(e)
	}
	portStr := ""
	if port != 0 {
		portStr = fmt.Sprintf(" port=%d", port)
	}
	_, err := fmt.Fprintf(l.w, "%s [%s]%s %s\n", e.Timestamp.Format(time.RFC3339), e.Kind, portStr, e.Message)
	return err
}
