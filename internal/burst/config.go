package burst

import (
	"errors"
	"time"
)

// Config holds parameters for the burst Tracker.
type Config struct {
	// Window is the sliding time period over which events are counted.
	Window time.Duration
	// Ceiling is the maximum number of events allowed before Spiking returns true.
	Ceiling int
}

// DefaultConfig returns a sensible default configuration.
func DefaultConfig() Config {
	return Config{
		Window:  30 * time.Second,
		Ceiling: 10,
	}
}

// Validate returns an error if the Config contains invalid values.
func (c Config) Validate() error {
	if c.Window <= 0 {
		return errors.New("burst: window must be positive")
	}
	if c.Ceiling < 0 {
		return errors.New("burst: ceiling must be non-negative")
	}
	return nil
}

// NewFromConfig validates cfg and returns a ready-to-use Tracker.
func NewFromConfig(cfg Config) (*Tracker, error) {
	return New(cfg)
}
