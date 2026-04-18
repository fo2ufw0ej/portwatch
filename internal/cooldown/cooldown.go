// Package cooldown provides a per-key cooldown tracker that prevents
// repeated actions within a configurable quiet period.
package cooldown

import (
	"sync"
	"time"
)

// Cooldown tracks the last allowed time per key.
type Cooldown struct {
	mu      sync.Mutex
	period  time.Duration
	entries map[string]time.Time
	now     func() time.Time
}

// New returns a Cooldown with the given quiet period.
func New(period time.Duration) *Cooldown {
	return &Cooldown{
		period:  period,
		entries: make(map[string]time.Time),
		now:     time.Now,
	}
}

// Allow returns true if the key has not been seen within the cooldown period.
// If allowed, the key's timestamp is updated.
func (c *Cooldown) Allow(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.now()
	if last, ok := c.entries[key]; ok {
		if now.Sub(last) < c.period {
			return false
		}
	}
	c.entries[key] = now
	return true
}

// Reset clears the cooldown entry for the given key.
func (c *Cooldown) Reset(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Remaining returns how long until the key's cooldown expires.
// Returns 0 if the key is not in cooldown.
func (c *Cooldown) Remaining(key string) time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()
	last, ok := c.entries[key]
	if !ok {
		return 0
	}
	remaining := c.period - c.now().Sub(last)
	if remaining < 0 {
		return 0
	}
	return remaining
}
