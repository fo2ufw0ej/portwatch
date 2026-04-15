// Package ratelimit provides a simple token-bucket style rate limiter
// to suppress alert flooding when ports change frequently.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter tracks alert events per key and suppresses bursts.
type Limiter struct {
	mu       sync.Mutex
	window   time.Duration
	maxBurst int
	buckets  map[string][]time.Time
	now      func() time.Time
}

// New creates a Limiter that allows at most maxBurst events per key
// within the given window duration.
func New(window time.Duration, maxBurst int) *Limiter {
	return &Limiter{
		window:   window,
		maxBurst: maxBurst,
		buckets:  make(map[string][]time.Time),
		now:      time.Now,
	}
}

// Allow returns true if the event identified by key is within the
// allowed burst limit for the current window.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	cutoff := now.Add(-l.window)

	times := l.buckets[key]
	filtered := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}

	if len(filtered) >= l.maxBurst {
		l.buckets[key] = filtered
		return false
	}

	l.buckets[key] = append(filtered, now)
	return true
}

// Reset clears all recorded events for the given key.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.buckets, key)
}

// Len returns the number of events recorded for key within the current window.
func (l *Limiter) Len(key string) int {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	cutoff := now.Add(-l.window)
	count := 0
	for _, t := range l.buckets[key] {
		if t.After(cutoff) {
			count++
		}
	}
	return count
}
