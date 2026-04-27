package latch

import "fmt"

// Config holds optional configuration for a Latch.
// Currently it is a placeholder for future options such as
// automatic expiry of tripped entries.
type Config struct {
	// MaxEntries caps the number of simultaneously tripped ports.
	// Zero means unlimited.
	MaxEntries int
}

// DefaultConfig returns a Config with permissive defaults.
func DefaultConfig() Config {
	return Config{
		MaxEntries: 0,
	}
}

// Validate returns an error if the Config is invalid.
func (c Config) Validate() error {
	if c.MaxEntries < 0 {
		return fmt.Errorf("latch: MaxEntries must be >= 0, got %d", c.MaxEntries)
	}
	return nil
}

// NewFromConfig returns a Latch built from cfg.
func NewFromConfig(cfg Config) (*Latch, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return New(), nil
}
