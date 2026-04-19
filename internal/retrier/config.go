package retrier

import (
	"errors"
	"time"
)

// Config holds parameters for retry behaviour.
type Config struct {
	MaxAttempts     int
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
}

// DefaultConfig returns a sensible default retry config.
func DefaultConfig() Config {
	return Config{
		MaxAttempts:     3,
		InitialInterval: 200 * time.Millisecond,
		MaxInterval:     5 * time.Second,
		Multiplier:      2.0,
	}
}

// Validate checks that all config fields are valid.
func (c Config) Validate() error {
	if c.MaxAttempts < 1 {
		return errors.New("retrier: MaxAttempts must be at least 1")
	}
	if c.InitialInterval <= 0 {
		return errors.New("retrier: InitialInterval must be positive")
	}
	if c.MaxInterval < c.InitialInterval {
		return errors.New("retrier: MaxInterval must be >= InitialInterval")
	}
	if c.Multiplier < 1.0 {
		return errors.New("retrier: Multiplier must be >= 1.0")
	}
	return nil
}

// NewFromConfig validates and constructs a Retrier from config.
func NewFromConfig(cfg Config) (*Retrier, error) {
	return New(cfg)
}

// WithMaxAttempts returns a copy of the config with MaxAttempts set to n.
func (c Config) WithMaxAttempts(n int) Config {
	c.MaxAttempts = n
	return c
}

// WithInitialInterval returns a copy of the config with InitialInterval set to d.
func (c Config) WithInitialInterval(d time.Duration) Config {
	c.InitialInterval = d
	return c
}
