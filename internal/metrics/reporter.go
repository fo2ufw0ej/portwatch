package metrics

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// Format controls output style for metrics reporting.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Write serialises a Snapshot to w in the given format.
// If w is nil it falls back to os.Stdout.
func Write(w io.Writer, s Snapshot, f Format) error {
	if w == nil {
		w = os.Stdout
	}
	switch f {
	case FormatJSON:
		return writeJSON(w, s)
	default:
		return writeText(w, s)
	}
}

func writeText(w io.Writer, s Snapshot) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Uptime:\t%s\n", s.Uptime)
	fmt.Fprintf(tw, "Scans total:\t%d\n", s.ScansTotal)
	fmt.Fprintf(tw, "Alerts total:\t%d\n", s.AlertsTotal)
	if !s.LastScanAt.IsZero() {
		fmt.Fprintf(tw, "Last scan at:\t%s\n", s.LastScanAt.Format("2006-01-02 15:04:05"))
		fmt.Fprintf(tw, "Last scan ports:\t%d\n", s.LastScanPorts)
	}
	return tw.Flush()
}

func writeJSON(w io.Writer, s Snapshot) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}
