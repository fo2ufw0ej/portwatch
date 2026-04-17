package backoff

import (
	"errors"
	"time"
)

// Validate checks that cfg fields are sensible.
func (c Config) Validate() error {
	if c.InitialInterval <= 0 {
		return errors.New("backoff: InitialInterval must be positive")
	}
	if c.MaxInterval < c.InitialInterval {
		return errors.New("backoff: MaxInterval must be >= InitialInterval")
	}
	if c.Multiplier < 1.0 {
		return errors.New("backoff: Multiplier must be >= 1.0")
	}
	if c.Jitter < 0 || c.Jitter > 1 {
		return errors.New("backoff: Jitter must be in [0, 1]")
	}
	return nil
}

// NewFromConfig validates cfg and returns a ready Backoff.
func NewFromConfig(cfg Config) (*Backoff, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return New(cfg), nil
}

// NewDefault returns a Backoff using DefaultConfig.
func NewDefault() *Backoff {
	return New(DefaultConfig())
}

// WithMaxInterval returns a copy of cfg with MaxInterval set to d.
func (c Config) WithMaxInterval(d time.Duration) Config {
	c.MaxInterval = d
	return c
}
