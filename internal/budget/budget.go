// Package budget implements an error budget tracker that measures the
// fraction of successful operations over a sliding window. When the budget
// is exhausted (success rate drops below the configured threshold) callers
// can back off or suppress non-critical alerts.
package budget

import (
	"fmt"
	"sync"
	"time"
)

// Budget tracks successes and failures within a rolling time window and
// exposes the remaining error budget as a value in [0, 1].
type Budget struct {
	mu       sync.Mutex
	cfg      Config
	buckets  []bucket
	current  int
	lastTick time.Time
}

type bucket struct {
	successes int
	failures  int
}

// New creates a Budget from cfg. Returns an error if cfg is invalid.
func New(cfg Config) (*Budget, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("budget: %w", err)
	}
	return &Budget{
		cfg:      cfg,
		buckets:  make([]bucket, cfg.Buckets),
		lastTick: time.Now(),
	}, nil
}

// Record registers one observation. success=true counts as a success,
// false as a failure.
func (b *Budget) Record(success bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.advance()
	if success {
		b.buckets[b.current].successes++
	} else {
		b.buckets[b.current].failures++
	}
}

// Remaining returns the fraction of error budget still available in [0, 1].
// A value of 1.0 means no errors; 0.0 means fully exhausted.
func (b *Budget) Remaining() float64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.advance()
	var s, f int
	for _, bk := range b.buckets {
		s += bk.successes
		f += bk.failures
	}
	total := s + f
	if total == 0 {
		return 1.0
	}
	actualRate := float64(s) / float64(total)
	remaining := (actualRate - b.cfg.Threshold) / (1.0 - b.cfg.Threshold)
	if remaining < 0 {
		return 0
	}
	if remaining > 1 {
		return 1
	}
	return remaining
}

// Exhausted reports whether the error budget has been fully consumed.
func (b *Budget) Exhausted() bool { return b.Remaining() == 0 }

// Reset clears all recorded observations.
func (b *Budget) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i := range b.buckets {
		b.buckets[i] = bucket{}
	}
}

// advance rotates expired buckets.
func (b *Budget) advance() {
	now := time.Now()
	bucketDur := b.cfg.Window / time.Duration(b.cfg.Buckets)
	for now.Sub(b.lastTick) >= bucketDur {
		b.current = (b.current + 1) % b.cfg.Buckets
		b.buckets[b.current] = bucket{}
		b.lastTick = b.lastTick.Add(bucketDur)
	}
}
