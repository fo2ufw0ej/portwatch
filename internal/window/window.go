// Package window provides a sliding time-window counter for tracking
// event frequency over a rolling duration.
package window

import (
	"sync"
	"time"
)

// Counter tracks event counts within a sliding time window.
type Counter struct {
	mu       sync.Mutex
	period   time.Duration
	buckets  []bucket
	size     int
	current  int
}

type bucket struct {
	at    time.Time
	count int
}

// New creates a Counter with the given window period and bucket count.
func New(period time.Duration, buckets int) *Counter {
	if buckets < 1 {
		buckets = 1
	}
	return &Counter{
		period:  period,
		buckets: make([]bucket, buckets),
		size:    buckets,
	}
}

// Add records n events at the current time.
func (c *Counter) Add(n int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	c.buckets[c.current] = bucket{at: now, count: c.buckets[c.current].count + n}
}

// Tick advances to the next bucket, resetting its count.
func (c *Counter) Tick() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.current = (c.current + 1) % c.size
	c.buckets[c.current] = bucket{}
}

// Total returns the sum of events within the active window period.
func (c *Counter) Total() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	cutoff := time.Now().Add(-c.period)
	total := 0
	for _, b := range c.buckets {
		if b.at.After(cutoff) {
			total += b.count
		}
	}
	return total
}

// Reset clears all buckets.
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i := range c.buckets {
		c.buckets[i] = bucket{}
	}
	c.current = 0
}
