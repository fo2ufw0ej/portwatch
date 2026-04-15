// Package suppress provides a mechanism to suppress duplicate alerts
// for ports that have already been reported within a cooldown window.
package suppress

import (
	"sync"
	"time"
)

// Entry tracks the last alert time and suppression state for a port.
type Entry struct {
	LastAlert time.Time
	Count     int
}

// Suppressor tracks which port events have recently been alerted
// and suppresses duplicates within a configurable cooldown window.
type Suppressor struct {
	mu       sync.Mutex
	cooldown time.Duration
	entries  map[string]*Entry
	now      func() time.Time
}

// New creates a new Suppressor with the given cooldown duration.
func New(cooldown time.Duration) *Suppressor {
	return &Suppressor{
		cooldown: cooldown,
		entries:  make(map[string]*Entry),
		now:      time.Now,
	}
}

// Allow returns true if the given key should trigger an alert,
// and false if it is within the cooldown window of a previous alert.
// If allowed, the entry is updated with the current time.
func (s *Suppressor) Allow(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.now()
	e, ok := s.entries[key]
	if !ok || now.Sub(e.LastAlert) >= s.cooldown {
		s.entries[key] = &Entry{LastAlert: now, Count: 1}
		if ok {
			s.entries[key].Count = e.Count + 1
		}
		return true
	}
	return false
}

// Reset clears suppression state for the given key.
func (s *Suppressor) Reset(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, key)
}

// Stats returns a snapshot of the current suppression entries.
func (s *Suppressor) Stats() map[string]Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make(map[string]Entry, len(s.entries))
	for k, v := range s.entries {
		out[k] = *v
	}
	return out
}
