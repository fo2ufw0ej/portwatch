package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/yourorg/portwatch/internal/history"
)

// WriteHistory formats and writes history entries to w.
// format must be "text" or "json".
func WriteHistory(w io.Writer, entries []history.Entry, format string) error {
	if w == nil {
		w = os.Stdout
	}
	if len(entries) == 0 {
		_, err := fmt.Fprintln(w, "No history recorded.")
		return err
	}

	switch format {
	case "json":
		return writeHistoryJSON(w, entries)
	default:
		return writeHistoryText(w, entries)
	}
}

func writeHistoryText(w io.Writer, entries []history.Entry) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIMESTAMP\tOPENED\tCLOSED")
	for _, e := range entries {
		ts := e.Timestamp.Format(time.RFC3339)
		opened := PortListString(e.Opened)
		closed := PortListString(e.Closed)
		if opened == "" {
			opened = "-"
		}
		if closed == "" {
			closed = "-"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\n", ts, opened, closed)
	}
	return tw.Flush()
}

func writeHistoryJSON(w io.Writer, entries []history.Entry) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}
