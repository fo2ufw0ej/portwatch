// Package stride tracks the rate of port change events over time,
// providing a rolling events-per-second measurement useful for detecting
// scanning activity or rapid churn.
package stride

import (
	"sync"
	"time"
)

// Stride tracks an event rate using a sliding window of timestamps.
type Stride struct {
	mu       sync.Mutex
	window   time.Duration
	times    []time.Time
	nowFn    func() time.Time
}

// New returns a Stride that measures rate over the given window duration.
func New(window time.Duration) *Stride {
	if window <= 0 {
		window = 10 * time.Second
	}
	return &Stride{
		window: window,
		nowFn:  time.Now,
	}
}

// Record registers a new event at the current time.
func (s *Stride) Record() {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.nowFn()
	s.prune(now)
	s.times = append(s.times, now)
}

// Rate returns the number of events per second observed within the window.
func (s *Stride) Rate() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.nowFn()
	s.prune(now)
	if len(s.times) == 0 {
		return 0
	}
	return float64(len(s.times)) / s.window.Seconds()
}

// Count returns the raw number of events within the current window.
func (s *Stride) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prune(s.nowFn())
	return len(s.times)
}

// Reset clears all recorded events.
func (s *Stride) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.times = s.times[:0]
}

// prune removes events older than the window. Must be called with mu held.
func (s *Stride) prune(now time.Time) {
	cutoff := now.Add(-s.window)
	i := 0
	for i < len(s.times) && s.times[i].Before(cutoff) {
		i++
	}
	s.times = s.times[i:]
}
