package grace

import (
	"errors"
	"time"
)

// Config holds configuration for the graceful shutdown Coordinator.
type Config struct {
	// Timeout is the maximum time to wait for components to finish.
	Timeout time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Timeout: 10 * time.Second,
	}
}

// Validate checks that the Config is valid.
func (c Config) Validate() error {
	if c.Timeout <= 0 {
		return errors.New("grace: timeout must be positive")
	}
	return nil
}

// NewFromConfig returns a Coordinator from the given Config.
func NewFromConfig(cfg Config) (*Coordinator, error) {
	return New(cfg)
}
