package budget

import (
	"errors"
	"fmt"
	"time"
)

// Config holds parameters for a Budget.
type Config struct {
	// Window is the total duration of the sliding window.
	Window time.Duration
	// Buckets is the number of sub-windows used for resolution.
	Buckets int
	// Threshold is the minimum acceptable success rate in [0, 1).
	// When the actual success rate falls below this value the budget is
	// considered exhausted.
	Threshold float64
}

// DefaultConfig returns a sensible default configuration.
func DefaultConfig() Config {
	return Config{
		Window:    5 * time.Minute,
		Buckets:   10,
		Threshold: 0.95,
	}
}

// Validate returns an error if the configuration is invalid.
func (c Config) Validate() error {
	if c.Window <= 0 {
		return errors.New("window must be positive")
	}
	if c.Buckets < 1 {
		return errors.New("buckets must be at least 1")
	}
	if c.Threshold < 0 || c.Threshold >= 1 {
		return fmt.Errorf("threshold must be in [0, 1), got %v", c.Threshold)
	}
	return nil
}

// NewFromConfig is an alias for New that makes the construction pattern
// consistent with other packages in this module.
func NewFromConfig(cfg Config) (*Budget, error) { return New(cfg) }
