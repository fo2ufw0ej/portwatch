// Package shadow maintains a shadow copy of the last-known port state
// and detects whether the current scan differs from it.
package shadow

import (
	"sync"
	"time"
)

// Entry holds a snapshot of ports at a point in time.
type Entry struct {
	Ports     []int
	RecordedAt time.Time
}

// Store keeps the most recent port snapshot in memory.
type Store struct {
	mu      sync.RWMutex
	current *Entry
}

// New returns an empty shadow Store.
func New() *Store {
	return &Store{}
}

// Set replaces the stored entry with the given ports.
func (s *Store) Set(ports []int) {
	copied := make([]int, len(ports))
	copy(copied, ports)
	s.mu.Lock()
	s.current = &Entry{Ports: copied, RecordedAt: time.Now()}
	s.mu.Unlock()
}

// Get returns the current shadow entry, or nil if none has been set.
func (s *Store) Get() *Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.current == nil {
		return nil
	}
	copied := make([]int, len(s.current.Ports))
	copy(copied, s.current.Ports)
	return &Entry{Ports: copied, RecordedAt: s.current.RecordedAt}
}

// Changed reports whether ports differs from the stored shadow.
// If no shadow exists, Changed returns true so the caller always
// persists the first observation.
func (s *Store) Changed(ports []int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.current == nil {
		return true
	}
	if len(ports) != len(s.current.Ports) {
		return true
	}
	set := make(map[int]struct{}, len(s.current.Ports))
	for _, p := range s.current.Ports {
		set[p] = struct{}{}
	}
	for _, p := range ports {
		if _, ok := set[p]; !ok {
			return true
		}
	}
	return false
}

// Clear removes the stored entry.
func (s *Store) Clear() {
	s.mu.Lock()
	s.current = nil
	s.mu.Unlock()
}
