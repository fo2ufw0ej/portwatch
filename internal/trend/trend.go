// Package trend tracks port change frequency over time, producing a
// simple rising/falling/stable signal useful for suppressing noise.
package trend

import (
	"sync"
	"time"
)

// Direction describes the current direction of change activity.
type Direction string

const (
	Stable  Direction = "stable"
	Rising  Direction = "rising"
	Falling Direction = "falling"
)

// Sample is a single observation recorded at a point in time.
type Sample struct {
	At    time.Time
	Delta int // positive = ports opened, negative = ports closed
}

// Tracker accumulates port-change samples and derives a trend direction.
type Tracker struct {
	mu      sync.Mutex
	window  time.Duration
	samples []Sample
}

// New returns a Tracker that considers samples within the given window.
func New(window time.Duration) *Tracker {
	if window <= 0 {
		window = time.Minute
	}
	return &Tracker{window: window}
}

// Record adds a sample representing net port changes (opened minus closed).
func (t *Tracker) Record(delta int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.prune(time.Now())
	t.samples = append(t.samples, Sample{At: time.Now(), Delta: delta})
}

// Direction returns the current trend based on samples within the window.
func (t *Tracker) Direction() Direction {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.prune(time.Now())

	var sum int
	for _, s := range t.samples {
		sum += s.Delta
	}
	switch {
	case sum > 0:
		return Rising
	case sum < 0:
		return Falling
	default:
		return Stable
	}
}

// Samples returns a copy of the samples currently within the window.
func (t *Tracker) Samples() []Sample {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.prune(time.Now())
	out := make([]Sample, len(t.samples))
	copy(out, t.samples)
	return out
}

// Reset discards all recorded samples.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.samples = nil
}

// prune removes samples older than the window. Caller must hold t.mu.
func (t *Tracker) prune(now time.Time) {
	cutoff := now.Add(-t.window)
	i := 0
	for i < len(t.samples) && t.samples[i].At.Before(cutoff) {
		i++
	}
	t.samples = t.samples[i:]
}
