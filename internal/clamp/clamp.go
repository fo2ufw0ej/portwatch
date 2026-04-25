// Package clamp provides a port-count clamping utility that enforces
// configurable minimum and maximum bounds on observed port counts,
// emitting a warning when a value is clamped.
package clamp

import "fmt"

// Config holds the bounds for the clamper.
type Config struct {
	Min int
	Max int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Min: 0,
		Max: 65535,
	}
}

// Validate returns an error if the config is invalid.
func (c Config) Validate() error {
	if c.Min < 0 {
		return fmt.Errorf("clamp: min must be >= 0, got %d", c.Min)
	}
	if c.Max < c.Min {
		return fmt.Errorf("clamp: max (%d) must be >= min (%d)", c.Max, c.Min)
	}
	return nil
}

// Clamper enforces min/max bounds on integer values.
type Clamper struct {
	cfg Config
}

// New returns a new Clamper or an error if the config is invalid.
func New(cfg Config) (*Clamper, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Clamper{cfg: cfg}, nil
}

// Clamp returns v clamped to [min, max] and a boolean indicating
// whether the value was changed.
func (c *Clamper) Clamp(v int) (int, bool) {
	if v < c.cfg.Min {
		return c.cfg.Min, true
	}
	if v > c.cfg.Max {
		return c.cfg.Max, true
	}
	return v, false
}

// ClampAll applies Clamp to every element of vs, returning the
// clamped slice and the count of values that were modified.
func (c *Clamper) ClampAll(vs []int) ([]int, int) {
	out := make([]int, len(vs))
	changed := 0
	for i, v := range vs {
		clamped, ok := c.Clamp(v)
		out[i] = clamped
		if ok {
			changed++
		}
	}
	return out, changed
}

// Bounds returns the configured [min, max] pair.
func (c *Clamper) Bounds() (int, int) {
	return c.cfg.Min, c.cfg.Max
}
