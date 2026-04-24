// Package burst tracks short-term spike activity over a sliding window.
// It counts events within a configurable period and reports whether the
// observed rate exceeds a defined ceiling.
package burst

import (
	"sync"
	"time"
)

// Tracker counts events per key and reports spikes.
type Tracker struct {
	mu      sync.Mutex
	cfg     Config
	buckets map[string][]time.Time
}

// New returns a Tracker using the provided Config.
// Returns an error if the config is invalid.
func New(cfg Config) (*Tracker, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Tracker{
		cfg:     cfg,
		buckets: make(map[string][]time.Time),
	}, nil
}

// Record registers one event for key at the current time.
func (t *Tracker) Record(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := time.Now()
	t.buckets[key] = append(t.prune(key, now), now)
}

// Count returns the number of events recorded for key within the window.
func (t *Tracker) Count(key string) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.prune(key, time.Now()))
}

// Spiking returns true when the event count for key exceeds the ceiling.
func (t *Tracker) Spiking(key string) bool {
	return t.Count(key) > t.cfg.Ceiling
}

// Reset clears all recorded events for key.
func (t *Tracker) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.buckets, key)
}

// prune removes timestamps outside the window and updates the bucket.
// Must be called with mu held.
func (t *Tracker) prune(key string, now time.Time) []time.Time {
	cutoff := now.Add(-t.cfg.Window)
	old := t.buckets[key]
	var fresh []time.Time
	for _, ts := range old {
		if ts.After(cutoff) {
			fresh = append(fresh, ts)
		}
	}
	t.buckets[key] = fresh
	return fresh
}
