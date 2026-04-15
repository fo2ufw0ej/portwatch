package reporter

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Format defines the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Reporter writes port scan summaries to an output destination.
type Reporter struct {
	out    io.Writer
	format Format
}

// New creates a Reporter writing to the given writer in the given format.
// If out is nil, os.Stdout is used.
func New(out io.Writer, format Format) *Reporter {
	if out == nil {
		out = os.Stdout
	}
	if format == "" {
		format = FormatText
	}
	return &Reporter{out: out, format: format}
}

// WriteSummary writes a human-readable or JSON summary of open ports.
func (r *Reporter) WriteSummary(ports []int) error {
	timestamp := time.Now().Format(time.RFC3339)
	switch r.format {
	case FormatJSON:
		return r.writeJSON(timestamp, ports)
	default:
		return r.writeText(timestamp, ports)
	}
}

func (r *Reporter) writeText(ts string, ports []int) error {
	_, err := fmt.Fprintf(r.out, "[%s] Open ports (%d): %v\n", ts, len(ports), ports)
	return err
}

func (r *Reporter) writeJSON(ts string, ports []int) error {
	// Marshal manually to avoid importing encoding/json for a trivial struct.
	portList := "["
	for i, p := range ports {
		if i > 0 {
			portList += ","
		}
		portList += fmt.Sprintf("%d", p)
	}
	portList += "]"
	_, err := fmt.Fprintf(r.out, `{"timestamp":%q,"open_port_count":%d,"open_ports":%s}\n`,
		ts, len(ports), portList)
	return err
}

// WriteDiff writes a summary of port changes.
func (r *Reporter) WriteDiff(diff scanner.Diff) error {
	timestamp := time.Now().Format(time.RFC3339)
	if len(diff.Opened) == 0 && len(diff.Closed) == 0 {
		return nil
	}
	_, err := fmt.Fprintf(r.out, "[%s] Changes — opened: %v  closed: %v\n",
		timestamp, diff.Opened, diff.Closed)
	return err
}
