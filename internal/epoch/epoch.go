// Package epoch tracks a monotonically increasing generation counter that
// increments each time a full scan cycle completes. Downstream components
// can use the epoch number to detect stale data or correlate events that
// belong to the same scan round.
package epoch

import (
	"fmt"
	"sync"
	"time"
)

// Epoch holds a single generation snapshot.
type Epoch struct {
	Number    uint64    `json:"number"`
	StartedAt time.Time `json:"started_at"`
}

// String returns a human-readable representation of the epoch.
func (e Epoch) String() string {
	return fmt.Sprintf("epoch#%d @ %s", e.Number, e.StartedAt.Format(time.RFC3339))
}

// Counter is a thread-safe, monotonically increasing epoch counter.
type Counter struct {
	mu      sync.RWMutex
	current Epoch
}

// New returns a Counter whose first epoch number is 0 (not yet advanced).
func New() *Counter {
	return &Counter{}
}

// Advance increments the epoch number and records the current time as the
// start of the new epoch. It returns the newly created Epoch.
func (c *Counter) Advance() Epoch {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.current = Epoch{
		Number:    c.current.Number + 1,
		StartedAt: time.Now(),
	}
	return c.current
}

// Current returns the most recent epoch without advancing it.
// Before the first call to Advance, Number is 0 and StartedAt is the zero
// time, which callers can use to detect an uninitialised counter.
func (c *Counter) Current() Epoch {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.current
}

// Reset sets the counter back to its zero state. Primarily useful in tests.
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.current = Epoch{}
}
