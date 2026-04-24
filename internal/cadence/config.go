package cadence

import (
	"errors"
	"time"
)

// Config holds parameters for a Tracker.
type Config struct {
	// Expected is the nominal interval between observations.
	Expected time.Duration

	// Tolerance is the maximum deviation from Expected that is still
	// considered "on time".
	Tolerance time.Duration

	// WindowSize is the maximum number of recent intervals to retain
	// for the Average calculation.
	WindowSize int
}

// DefaultConfig returns a Config suitable for a one-minute scan loop.
func DefaultConfig() Config {
	return Config{
		Expected:   60 * time.Second,
		Tolerance:  5 * time.Second,
		WindowSize: 10,
	}
}

// Validate returns an error if any field is out of range.
func (c Config) Validate() error {
	if c.Expected <= 0 {
		return errors.New("cadence: Expected must be positive")
	}
	if c.Tolerance < 0 {
		return errors.New("cadence: Tolerance must not be negative")
	}
	if c.WindowSize <= 0 {
		return errors.New("cadence: WindowSize must be positive")
	}
	return nil
}

// NewFromConfig is an alias for New that makes construction uniform
// with other packages in portwatch.
func NewFromConfig(cfg Config) (*Tracker, error) {
	return New(cfg)
}
