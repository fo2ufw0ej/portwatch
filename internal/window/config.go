package window

import (
	"fmt"
	"time"
)

// Config holds configuration for a sliding window counter.
type Config struct {
	Period  time.Duration `json:"period"`
	Buckets int           `json:"buckets"`
}

// DefaultConfig returns a sensible default window configuration.
func DefaultConfig() Config {
	return Config{
		Period:  time.Minute,
		Buckets: 6,
	}
}

// Validate checks that the config values are usable.
func (c Config) Validate() error {
	if c.Period <= 0 {
		return fmt.Errorf("window: period must be positive, got %s", c.Period)
	}
	if c.Buckets < 1 {
		return fmt.Errorf("window: buckets must be at least 1, got %d", c.Buckets)
	}
	return nil
}

// NewFromConfig constructs a Counter from a validated Config.
func NewFromConfig(cfg Config) (*Counter, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return New(cfg.Period, cfg.Buckets), nil
}
