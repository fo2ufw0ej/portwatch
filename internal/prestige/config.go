package prestige

import "fmt"

// Config holds configuration for the Tracker.
type Config struct {
	// Alpha is the EMA smoothing factor in (0, 1].
	Alpha float64 `json:"alpha" yaml:"alpha"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{Alpha: DefaultAlpha}
}

// Validate returns an error if the Config is invalid.
func (c Config) Validate() error {
	if c.Alpha <= 0 || c.Alpha > 1 {
		return fmt.Errorf("prestige: alpha must be in (0, 1], got %v", c.Alpha)
	}
	return nil
}

// NewFromConfig constructs a Tracker from a validated Config.
func NewFromConfig(c Config) (*Tracker, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return New(c.Alpha), nil
}
