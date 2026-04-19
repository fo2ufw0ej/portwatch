package decay

import "time"

// NewFromConfig creates a Scorer from the provided Config, returning an error
// if validation fails.
func NewFromConfig(cfg Config) (*Scorer, error) {
	return New(cfg)
}

// NewDefault creates a Scorer with default configuration.
func NewDefault() (*Scorer, error) {
	return New(DefaultConfig())
}

// WithHalfLife returns a copy of cfg with the given half-life.
func (c Config) WithHalfLife(d time.Duration) Config {
	c.HalfLife = d
	return c
}
