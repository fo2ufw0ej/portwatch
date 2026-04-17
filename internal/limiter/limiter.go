// Package limiter provides a scan-rate limiter that enforces a minimum
// interval between consecutive scans to avoid overwhelming the host.
package limiter

import (
	"errors"
	"sync"
	"time"
)

// Config holds limiter settings.
type Config struct {
	// MinInterval is the minimum duration required between scans.
	MinInterval time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		MinInterval: 5 * time.Second,
	}
}

// Validate returns an error if the config is invalid.
func (c Config) Validate() error {
	if c.MinInterval <= 0 {
		return errors.New("limiter: MinInterval must be positive")
	}
	return nil
}

// Limiter enforces a minimum interval between scan cycles.
type Limiter struct {
	cfg  Config
	mu   sync.Mutex
	last time.Time
}

// New creates a Limiter from the given Config.
func New(cfg Config) (*Limiter, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Limiter{cfg: cfg}, nil
}

// Allow returns true if enough time has elapsed since the last allowed call.
// If allowed, it records the current time as the last scan time.
func (l *Limiter) Allow() bool {
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.last.IsZero() || now.Sub(l.last) >= l.cfg.MinInterval {
		l.last = now
		return true
	}
	return false
}

// Reset clears the last scan time, allowing the next call to Allow to succeed.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.last = time.Time{}
}

// NextAllowed returns the time at which the next scan will be permitted.
func (l *Limiter) NextAllowed() time.Time {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.last.IsZero() {
		return time.Now()
	}
	return l.last.Add(l.cfg.MinInterval)
}
