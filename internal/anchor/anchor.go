// Package anchor tracks which ports are considered "anchored" — stable,
// long-lived ports that should not trigger alerts when seen repeatedly.
// A port becomes anchored after being observed consistently for a minimum
// number of consecutive scans.
package anchor

import (
	"sync"
	"time"
)

// Entry holds the observation state for a single port.
type Entry struct {
	FirstSeen  time.Time
	LastSeen   time.Time
	Streak     int
	Anchored   bool
}

// Store tracks port observation streaks and promotes stable ports to anchored.
type Store struct {
	mu        sync.Mutex
	entries   map[int]*Entry
	threshold int
	now       func() time.Time
}

// New creates a Store that anchors a port after threshold consecutive observations.
func New(threshold int) *Store {
	if threshold <= 0 {
		threshold = 3
	}
	return &Store{
		entries:   make(map[int]*Entry),
		threshold: threshold,
		now:       time.Now,
	}
}

// Observe records a port as seen during the current scan cycle.
// Returns true if the port transitioned to anchored on this call.
func (s *Store) Observe(port int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.now()
	e, ok := s.entries[port]
	if !ok {
		s.entries[port] = &Entry{FirstSeen: now, LastSeen: now, Streak: 1}
		return false
	}

	e.LastSeen = now
	e.Streak++

	if !e.Anchored && e.Streak >= s.threshold {
		e.Anchored = true
		return true
	}
	return false
}

// IsAnchored reports whether port has been promoted to anchored status.
func (s *Store) IsAnchored(port int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[port]
	return ok && e.Anchored
}

// Reset clears the streak and anchored state for a port.
func (s *Store) Reset(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, port)
}

// Get returns a copy of the entry for port, and whether it exists.
func (s *Store) Get(port int) (Entry, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[port]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}
