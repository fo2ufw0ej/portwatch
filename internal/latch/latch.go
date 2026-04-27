// Package latch provides a sticky boolean gate that trips on a condition
// and remains tripped until explicitly reset. Useful for tracking whether
// a port has ever been seen in an unexpected state within a monitoring epoch.
package latch

import (
	"sync"
	"time"
)

// Latch is a sticky flag per key. Once tripped it stays tripped until Reset.
type Latch struct {
	mu      sync.Mutex
	entries map[int]entry
}

type entry struct {
	tripped   bool
	trippedAt time.Time
}

// New returns a ready-to-use Latch.
func New() *Latch {
	return &Latch{entries: make(map[int]entry)}
}

// Trip marks the given port key as tripped. If already tripped, this is a no-op.
// Returns true if this call caused the trip (i.e. it was not already tripped).
func (l *Latch) Trip(port int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	e := l.entries[port]
	if e.tripped {
		return false
	}
	l.entries[port] = entry{tripped: true, trippedAt: time.Now()}
	return true
}

// Tripped returns true if the given port has been tripped and not yet reset.
func (l *Latch) Tripped(port int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.entries[port].tripped
}

// TrippedAt returns the time the latch was tripped and whether it is set.
func (l *Latch) TrippedAt(port int) (time.Time, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	e := l.entries[port]
	return e.trippedAt, e.tripped
}

// Reset clears the tripped state for the given port.
func (l *Latch) Reset(port int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.entries, port)
}

// ResetAll clears all tripped entries.
func (l *Latch) ResetAll() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = make(map[int]entry)
}

// Len returns the number of currently tripped ports.
func (l *Latch) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.entries)
}
