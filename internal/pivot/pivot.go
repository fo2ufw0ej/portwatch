// Package pivot tracks the most frequently seen port state transitions
// and identifies which ports are the most volatile over a sliding window.
package pivot

import (
	"sync"
	"time"
)

// Entry records transition counts for a single port.
type Entry struct {
	Port        int
	Transitions int
	LastAt      time.Time
}

// Tracker counts open/close transitions per port within a time window.
type Tracker struct {
	mu      sync.Mutex
	window  time.Duration
	entries map[int][]time.Time
	now     func() time.Time
}

// New returns a Tracker that retains events within the given window duration.
func New(window time.Duration) *Tracker {
	return &Tracker{
		window:  window,
		entries: make(map[int][]time.Time),
		now:     time.Now,
	}
}

// Record registers a transition event for the given port.
func (t *Tracker) Record(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	t.prune(port, now)
	t.entries[port] = append(t.entries[port], now)
}

// Count returns the number of transitions recorded for port within the window.
func (t *Tracker) Count(port int) int {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.prune(port, t.now())
	return len(t.entries[port])
}

// Top returns up to n ports with the highest transition counts.
func (t *Tracker) Top(n int) []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	var result []Entry
	for port := range t.entries {
		t.prune(port, now)
		ev := t.entries[port]
		if len(ev) == 0 {
			continue
		}
		result = append(result, Entry{
			Port:        port,
			Transitions: len(ev),
			LastAt:      ev[len(ev)-1],
		})
	}

	// simple insertion sort — port counts are small
	for i := 1; i < len(result); i++ {
		for j := i; j > 0 && result[j].Transitions > result[j-1].Transitions; j-- {
			result[j], result[j-1] = result[j-1], result[j]
		}
	}
	if n > 0 && len(result) > n {
		result = result[:n]
	}
	return result
}

// Reset clears all recorded events.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = make(map[int][]time.Time)
}

// prune removes events outside the window. Must be called with lock held.
func (t *Tracker) prune(port int, now time.Time) {
	cutoff := now.Add(-t.window)
	ev := t.entries[port]
	i := 0
	for i < len(ev) && ev[i].Before(cutoff) {
		i++
	}
	if i > 0 {
		t.entries[port] = ev[i:]
	}
}
