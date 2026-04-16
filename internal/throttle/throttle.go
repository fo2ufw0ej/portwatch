// Package throttle provides a scan throttle that enforces a minimum
// interval between consecutive scans to avoid CPU/network spikes.
package throttle

import (
	"sync"
	"time"
)

// Throttle enforces a minimum gap between allowed operations.
type Throttle struct {
	mu       sync.Mutex
	minGap   time.Duration
	lastSeen map[string]time.Time
	now      func() time.Time
}

// New returns a Throttle with the given minimum gap between calls per key.
func New(minGap time.Duration) *Throttle {
	return &Throttle{
		minGap:   minGap,
		lastSeen: make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow returns true if enough time has elapsed since the last allowed call
// for the given key. It updates the last-seen timestamp on success.
func (t *Throttle) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	if last, ok := t.lastSeen[key]; ok {
		if now.Sub(last) < t.minGap {
			return false
		}
	}
	t.lastSeen[key] = now
	return true
}

// Reset clears the last-seen record for a key.
func (t *Throttle) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.lastSeen, key)
}

// SetNow replaces the time source (for testing).
func (t *Throttle) SetNow(fn func() time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.now = fn
}
