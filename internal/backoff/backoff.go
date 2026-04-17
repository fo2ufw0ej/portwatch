// Package backoff provides exponential backoff with jitter for retry logic.
package backoff

import (
	"math"
	"math/rand"
	"time"
)

// Config holds backoff parameters.
type Config struct {
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
	Jitter          float64 // fraction in [0, 1]
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		InitialInterval: 500 * time.Millisecond,
		MaxInterval:     30 * time.Second,
		Multiplier:      2.0,
		Jitter:          0.2,
	}
}

// Backoff computes successive wait durations.
type Backoff struct {
	cfg     Config
	attempt int
}

// New creates a Backoff from cfg.
func New(cfg Config) *Backoff {
	return &Backoff{cfg: cfg}
}

// Next returns the next wait duration and increments the attempt counter.
func (b *Backoff) Next() time.Duration {
	base := float64(b.cfg.InitialInterval) * math.Pow(b.cfg.Multiplier, float64(b.attempt))
	if base > float64(b.cfg.MaxInterval) {
		base = float64(b.cfg.MaxInterval)
	}
	jitter := (rand.Float64()*2 - 1) * b.cfg.Jitter * base
	d := time.Duration(base + jitter)
	if d < 0 {
		d = 0
	}
	b.attempt++
	return d
}

// Reset restarts the backoff sequence.
func (b *Backoff) Reset() {
	b.attempt = 0
}

// Attempt returns the current attempt index.
func (b *Backoff) Attempt() int {
	return b.attempt
}
