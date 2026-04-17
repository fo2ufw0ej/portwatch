package rollup

import (
	"errors"
	"time"
)

// Config holds tunable parameters for the Rollup.
type Config struct {
	// Window is how long to wait after the last diff before flushing.
	Window time.Duration `json:"window" yaml:"window"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Window: 2 * time.Second,
	}
}

// Validate returns an error if the Config is invalid.
func (c Config) Validate() error {
	if c.Window <= 0 {
		return errors.New("rollup: window must be positive")
	}
	return nil
}

// NewFromConfig constructs a Rollup from cfg and the given handler.
func NewFromConfig(cfg Config, handle Handler) (*Rollup, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return New(cfg.Window, handle), nil
}
