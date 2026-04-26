package stale

import (
	"errors"
	"time"
)

// Config holds configuration for the stale Tracker.
type Config struct {
	// IdleAfter is the duration after which a port is considered stale.
	IdleAfter time.Duration `toml:"idle_after" json:"idle_after"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		IdleAfter: 5 * time.Minute,
	}
}

// Validate returns an error if the Config is invalid.
func (c Config) Validate() error {
	if c.IdleAfter <= 0 {
		return errors.New("stale: idle_after must be positive")
	}
	return nil
}

// NewFromConfig creates a Tracker from the provided Config.
func NewFromConfig(c Config) (*Tracker, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return New(c.IdleAfter), nil
}
