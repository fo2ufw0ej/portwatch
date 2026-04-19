// Package scorecard aggregates port-level health scores into a summary report.
package scorecard

import (
	"fmt"
	"io"
	"os"
	"sort"
	"text/tabwriter"
)

// Entry holds the health score for a single port.
type Entry struct {
	Port  int
	Score float64 // 0.0 (unstable) to 1.0 (stable)
	Label string
}

// Card is a collection of scored port entries.
type Card struct {
	entries map[int]Entry
}

// New returns an empty Card.
func New() *Card {
	return &Card{entries: make(map[int]Entry)}
}

// Set records or updates the score for a port.
func (c *Card) Set(port int, score float64, label string) {
	c.entries[port] = Entry{Port: port, Score: score, Label: label}
}

// Get returns the Entry for a port and whether it exists.
func (c *Card) Get(port int) (Entry, bool) {
	e, ok := c.entries[port]
	return e, ok
}

// Entries returns all entries sorted by port number.
func (c *Card) Entries() []Entry {
	out := make([]Entry, 0, len(c.entries))
	for _, e := range c.entries {
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Port < out[j].Port })
	return out
}

// Len returns the number of tracked ports.
func (c *Card) Len() int { return len(c.entries) }

// Write renders the scorecard as a text table to w (stdout if nil).
func (c *Card) Write(w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PORT\tLABEL\tSCORE")
	for _, e := range c.Entries() {
		fmt.Fprintf(tw, "%d\t%s\t%.2f\n", e.Port, e.Label, e.Score)
	}
	tw.Flush()
}
