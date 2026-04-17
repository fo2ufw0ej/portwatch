package circuitbreaker

import (
	"errors"
	"time"
)

// Config holds configuration for a Breaker.
type Config struct {
	MaxFailures int           `json:"max_failures" yaml:"max_failures"`
	ResetAfter  time.Duration `json:"reset_after"  yaml:"reset_after"`
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxFailures: 5,
		ResetAfter:  30 * time.Second,
	}
}

// Validate returns an error if the config is invalid.
func (c Config) Validate() error {
	if c.MaxFailures <= 0 {
		return errors.New("circuitbreaker: max_failures must be > 0")
	}
	if c.ResetAfter <= 0 {
		return errors.New("circuitbreaker: reset_after must be > 0")
	}
	return nil
}

// NewFromConfig constructs a Breaker from a Config.
func NewFromConfig(c Config) (*Breaker, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return New(c.MaxFailures, c.ResetAfter), nil
}
