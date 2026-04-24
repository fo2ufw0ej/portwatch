// Package roster maintains a live registry of known ports, tracking
// when each port was first seen and last confirmed open.
package roster

import (
	"sync"
	"time"
)

// Entry holds metadata about a single tracked port.
type Entry struct {
	Port        int
	FirstSeen   time.Time
	LastSeen    time.Time
	SeenCount   int
}

// Roster is a thread-safe registry of observed open ports.
type Roster struct {
	mu      sync.RWMutex
	entries map[int]*Entry
	now     func() time.Time
}

// New creates an empty Roster.
func New() *Roster {
	return &Roster{
		entries: make(map[int]*Entry),
		now:     time.Now,
	}
}

// Touch records that the given port was observed at the current time.
// If the port is new, a fresh Entry is created; otherwise LastSeen and
// SeenCount are updated.
func (r *Roster) Touch(port int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := r.now()
	if e, ok := r.entries[port]; ok {
		e.LastSeen = now
		e.SeenCount++
		return
	}
	r.entries[port] = &Entry{
		Port:      port,
		FirstSeen: now,
		LastSeen:  now,
		SeenCount: 1,
	}
}

// Remove deletes a port from the roster.
func (r *Roster) Remove(port int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, port)
}

// Get returns the Entry for a port and whether it exists.
func (r *Roster) Get(port int) (Entry, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.entries[port]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// Len returns the number of ports currently tracked.
func (r *Roster) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.entries)
}

// Ports returns a snapshot of all currently tracked port numbers.
func (r *Roster) Ports() []int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]int, 0, len(r.entries))
	for p := range r.entries {
		out = append(out, p)
	}
	return out
}
