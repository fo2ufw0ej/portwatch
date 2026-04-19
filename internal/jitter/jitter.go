// Package jitter provides utilities for adding randomised jitter to durations.
// It is used throughout portwatch to spread out concurrent operations and
// avoid thundering-herd effects when multiple goroutines fire simultaneously.
package jitter

import (
	"fmt"
	"math/rand"
	"time"
)

// Config holds jitter parameters.
type Config struct {
	// Factor is the maximum fraction of the base duration to add as jitter.
	// Must be in the range [0, 1].
	Factor float64
}

// DefaultConfig returns a Config with a modest default jitter factor.
func DefaultConfig() Config {
	return Config{Factor: 0.1}
}

// Validate returns an error if the Config is invalid.
func (c Config) Validate() error {
	if c.Factor < 0 || c.Factor > 1 {
		return fmt.Errorf("jitter: factor must be in [0, 1], got %v", c.Factor)
	}
	return nil
}

// Applier adds jitter to durations.
type Applier struct {
	cfg Config
	rng *rand.Rand
}

// New creates an Applier using the provided Config.
func New(cfg Config) (*Applier, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Applier{
		cfg: cfg,
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

// Apply returns d plus a random jitter up to Factor*d.
func (a *Applier) Apply(d time.Duration) time.Duration {
	if a.cfg.Factor == 0 || d <= 0 {
		return d
	}
	max := float64(d) * a.cfg.Factor
	offset := time.Duration(a.rng.Float64() * max)
	return d + offset
}

// ApplyNegative returns d minus a random jitter up to Factor*d.
func (a *Applier) ApplyNegative(d time.Duration) time.Duration {
	if a.cfg.Factor == 0 || d <= 0 {
		return d
	}
	max := float64(d) * a.cfg.Factor
	offset := time.Duration(a.rng.Float64() * max)
	if offset > d {
		return 0
	}
	return d - offset
}
