// Package stale tracks ports that have not been seen recently and
// marks them as stale after a configurable idle duration.
package stale

import (
	"sync"
	"time"
)

// Entry holds the last-seen timestamp and stale state for a port.
type Entry struct {
	LastSeen time.Time
	Stale    bool
}

// Tracker monitors port activity and marks ports stale after IdleAfter.
type Tracker struct {
	mu       sync.Mutex
	entries  map[int]Entry
	idleAfter time.Duration
	now      func() time.Time
}

// New creates a Tracker with the given idle duration.
func New(idleAfter time.Duration) *Tracker {
	return &Tracker{
		entries:   make(map[int]Entry),
		idleAfter: idleAfter,
		now:       time.Now,
	}
}

// Touch records that port was seen at the current time and clears its stale flag.
func (t *Tracker) Touch(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries[port] = Entry{LastSeen: t.now(), Stale: false}
}

// Sweep marks any port not touched within IdleAfter as stale.
// It returns the list of newly-stale ports.
func (t *Tracker) Sweep() []int {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.now()
	var staled []int
	for port, e := range t.entries {
		if !e.Stale && now.Sub(e.LastSeen) > t.idleAfter {
			e.Stale = true
			t.entries[port] = e
			staled = append(staled, port)
		}
	}
	return staled
}

// IsStale reports whether the given port is currently marked stale.
func (t *Tracker) IsStale(port int) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[port]
	return ok && e.Stale
}

// Remove deletes all tracking state for port.
func (t *Tracker) Remove(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, port)
}

// Len returns the number of ports currently tracked.
func (t *Tracker) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.entries)
}
