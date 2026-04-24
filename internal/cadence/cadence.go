// Package cadence tracks how regularly a repeating event fires and
// reports whether the observed interval has drifted outside an
// acceptable tolerance band.
package cadence

import (
	"sync"
	"time"
)

// Tracker measures the interval between successive observations and
// exposes whether the cadence is within the configured tolerance.
type Tracker struct {
	mu        sync.Mutex
	cfg       Config
	last      time.Time
	intervals []time.Duration
}

// New returns a Tracker validated against cfg.
func New(cfg Config) (*Tracker, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Tracker{cfg: cfg}, nil
}

// Observe records the current time as a new observation.
// It returns the measured interval since the previous call, or zero
// if this is the first observation.
func (t *Tracker) Observe(now time.Time) time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.last.IsZero() {
		t.last = now
		return 0
	}

	interval := now.Sub(t.last)
	t.last = now

	t.intervals = append(t.intervals, interval)
	if len(t.intervals) > t.cfg.WindowSize {
		t.intervals = t.intervals[len(t.intervals)-t.cfg.WindowSize:]
	}
	return interval
}

// OnTime reports whether the most recent interval falls within
// [Expected - Tolerance, Expected + Tolerance].
// Returns false when fewer than two observations have been recorded.
func (t *Tracker) OnTime() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(t.intervals) == 0 {
		return false
	}
	last := t.intervals[len(t.intervals)-1]
	lo := t.cfg.Expected - t.cfg.Tolerance
	hi := t.cfg.Expected + t.cfg.Tolerance
	return last >= lo && last <= hi
}

// Average returns the mean interval over the observation window.
// Returns zero when no intervals have been recorded.
func (t *Tracker) Average() time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(t.intervals) == 0 {
		return 0
	}
	var sum time.Duration
	for _, iv := range t.intervals {
		sum += iv
	}
	return sum / time.Duration(len(t.intervals))
}

// Reset clears all recorded observations.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.last = time.Time{}
	t.intervals = t.intervals[:0]
}
