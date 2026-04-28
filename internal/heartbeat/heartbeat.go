// Package heartbeat provides a periodic liveness signal that records
// timestamps and exposes whether the last beat is within an expected window.
package heartbeat

import (
	"sync"
	"time"
)

// Heartbeat tracks the last time a beat was recorded and whether the
// signal is considered alive within a configurable TTL.
type Heartbeat struct {
	mu     sync.RWMutex
	lastAt time.Time
	ttl    time.Duration
	now    func() time.Time
}

// New returns a Heartbeat with the given time-to-live. A beat must be
// recorded within ttl for Alive to return true.
func New(ttl time.Duration) *Heartbeat {
	if ttl <= 0 {
		ttl = 30 * time.Second
	}
	return &Heartbeat{
		ttl: ttl,
		now: time.Now,
	}
}

// Beat records the current time as the most recent heartbeat.
func (h *Heartbeat) Beat() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.lastAt = h.now()
}

// LastAt returns the time of the most recent beat, or the zero value if
// Beat has never been called.
func (h *Heartbeat) LastAt() time.Time {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.lastAt
}

// Alive reports whether a beat was recorded within the configured TTL.
func (h *Heartbeat) Alive() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.lastAt.IsZero() {
		return false
	}
	return h.now().Sub(h.lastAt) <= h.ttl
}

// StaleSince returns how long ago the last beat occurred, or zero if
// the heartbeat is still considered alive.
func (h *Heartbeat) StaleSince() time.Duration {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.lastAt.IsZero() {
		return 0
	}
	elapsed := h.now().Sub(h.lastAt)
	if elapsed <= h.ttl {
		return 0
	}
	return elapsed - h.ttl
}

// Reset clears the last beat timestamp, causing Alive to return false
// until Beat is called again.
func (h *Heartbeat) Reset() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.lastAt = time.Time{}
}
