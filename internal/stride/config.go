package stride

import (
	"errors"
	"time"
)

// Config holds configuration for a Stride rate tracker.
type Config struct {
	// Window is the rolling duration over which events are counted.
	Window time.Duration `json:"window"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Window: 10 * time.Second,
	}
}

// Validate returns an error if the Config is invalid.
func (c Config) Validate() error {
	if c.Window <= 0 {
		return errors.New("stride: window must be positive")
	}
	return nil
}

// NewFromConfig creates a Stride from the given Config.
func NewFromConfig(c Config) (*Stride, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return New(c.Window), nil
}
