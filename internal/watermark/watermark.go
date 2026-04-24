// Package watermark tracks high-water and low-water marks for open port counts
// over time, enabling detection of unusual spikes or drops in port activity.
package watermark

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Mark holds the peak and trough port counts observed.
type Mark struct {
	High      int       `json:"high"`
	Low       int       `json:"low"`
	HighAt    time.Time `json:"high_at"`
	LowAt     time.Time `json:"low_at"`
	Observed  int       `json:"observed"`
}

// Watermark tracks high and low port count observations.
type Watermark struct {
	mu   sync.Mutex
	mark Mark
	init bool
}

// New returns a new Watermark with no observations recorded.
func New() *Watermark {
	return &Watermark{}
}

// Observe records a new port count sample, updating marks as needed.
func (w *Watermark) Observe(count int) {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now()
	w.mark.Observed++

	if !w.init {
		w.mark.High = count
		w.mark.Low = count
		w.mark.HighAt = now
		w.mark.LowAt = now
		w.init = true
		return
	}

	if count > w.mark.High {
		w.mark.High = count
		w.mark.HighAt = now
	}
	if count < w.mark.Low {
		w.mark.Low = count
		w.mark.LowAt = now
	}
}

// Snapshot returns a copy of the current marks.
func (w *Watermark) Snapshot() Mark {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.mark
}

// Reset clears all recorded marks.
func (w *Watermark) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.mark = Mark{}
	w.init = false
}

// Write prints the current watermark state to w; defaults to stdout if nil.
func (w *Watermark) Write(out io.Writer) {
	if out == nil {
		out = os.Stdout
	}
	w.mu.Lock()
	m := w.mark
	w.mu.Unlock()

	if !w.init {
		fmt.Fprintln(out, "watermark: no observations recorded")
		return
	}
	fmt.Fprintf(out, "watermark: high=%d (at %s)  low=%d (at %s)  observations=%d\n",
		m.High, m.HighAt.Format(time.RFC3339),
		m.Low, m.LowAt.Format(time.RFC3339),
		m.Observed)
}
