package trend

import (
	"fmt"
	"time"
)

// Config holds tunable parameters for a Tracker.
type Config struct {
	// Window is the rolling duration over which samples are retained.
	Window time.Duration `json:"window" yaml:"window"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Window: 5 * time.Minute,
	}
}

// Validate returns an error if the Config is invalid.
func (c Config) Validate() error {
	if c.Window <= 0 {
		return fmt.Errorf("trend: window must be positive, got %s", c.Window)
	}
	return nil
}

// NewFromConfig constructs a Tracker from a validated Config.
func NewFromConfig(c Config) (*Tracker, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return New(c.Window), nil
}
