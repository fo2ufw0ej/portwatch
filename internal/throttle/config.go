package throttle

import (
	"errors"
	"time"
)

// Config holds configuration for a Throttle instance.
type Config struct {
	// MinGap is the minimum duration that must elapse between allowed calls
	// for the same key.
	MinGap time.Duration `toml:"min_gap" json:"min_gap"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		MinGap: 5 * time.Second,
	}
}

// Validate returns an error if the config is invalid.
func (c Config) Validate() error {
	if c.MinGap <= 0 {
		return errors.New("throttle: min_gap must be positive")
	}
	return nil
}

// NewFromConfig constructs a Throttle from a Config after validating it.
func NewFromConfig(cfg Config) (*Throttle, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return New(cfg.MinGap), nil
}
