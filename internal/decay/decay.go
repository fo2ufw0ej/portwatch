// Package decay implements an exponential decay scorer that reduces
// a value toward zero over time, useful for aging out stale observations.
package decay

import (
	"math"
	"sync"
	"time"
)

// Config holds decay parameters.
type Config struct {
	HalfLife time.Duration // time for a value to halve
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{HalfLife: 5 * time.Minute}
}

// Validate returns an error if the config is invalid.
func (c Config) Validate() error {
	if c.HalfLife <= 0 {
		return ErrInvalidHalfLife
	}
	return nil
}

// entry holds a value and the last update time.
type entry struct {
	value     float64
	updatedAt time.Time
}

// Scorer tracks per-key decayed scores.
type Scorer struct {
	mu      sync.Mutex
	cfg     Config
	entries map[string]entry
}

// New creates a new Scorer with the given config.
func New(cfg Config) (*Scorer, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Scorer{cfg: cfg, entries: make(map[string]entry)}, nil
}

// Add adds delta to the decayed current value for key.
func (s *Scorer) Add(key string, delta float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	e := s.entries[key]
	decayed := s.decayValue(e.value, e.updatedAt, now)
	s.entries[key] = entry{value: decayed + delta, updatedAt: now}
}

// Get returns the current decayed score for key.
func (s *Scorer) Get(key string) float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	e := s.entries[key]
	return s.decayValue(e.value, e.updatedAt, time.Now())
}

// Reset clears the score for key.
func (s *Scorer) Reset(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, key)
}

func (s *Scorer) decayValue(v float64, since, now time.Time) float64 {
	if since.IsZero() || v == 0 {
		return v
	}
	elapsed := now.Sub(since).Seconds()
	hl := s.cfg.HalfLife.Seconds()
	return v * math.Pow(0.5, elapsed/hl)
}
