// Package sampler provides periodic port scan sampling with jitter to avoid
// thundering-herd effects when multiple instances run simultaneously.
package sampler

import (
	"context"
	"math/rand"
	"time"
)

// Config holds sampler configuration.
type Config struct {
	// Base interval between samples.
	Interval time.Duration
	// Jitter is the maximum random duration added to each interval.
	Jitter time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Interval: 30 * time.Second,
		Jitter:   5 * time.Second,
	}
}

// Validate returns an error string if the config is invalid.
func (c Config) Validate() error {
	if c.Interval <= 0 {
		return errInvalidInterval
	}
	if c.Jitter < 0 {
		return errNegativeJitter
	}
	return nil
}

var (
	errInvalidInterval = configError("interval must be positive")
	errNegativeJitter  = configError("jitter must be non-negative")
)

type configError string

func (e configError) Error() string { return string(e) }

// Sampler fires a callback at jittered intervals.
type Sampler struct {
	cfg Config
	rng *rand.Rand
}

// New creates a new Sampler.
func New(cfg Config) (*Sampler, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Sampler{
		cfg: cfg,
		//nolint:gosec — jitter does not require cryptographic randomness
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

// Run calls fn each interval (plus jitter) until ctx is cancelled.
func (s *Sampler) Run(ctx context.Context, fn func()) {
	for {
		wait := s.cfg.Interval
		if s.cfg.Jitter > 0 {
			wait += time.Duration(s.rng.Int63n(int64(s.cfg.Jitter)))
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(wait):
			fn()
		}
	}
}
