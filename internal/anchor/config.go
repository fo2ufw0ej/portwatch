package anchor

import "fmt"

// Config holds tunable parameters for the anchor store.
type Config struct {
	// Threshold is the number of consecutive observations required before a
	// port is considered anchored. Must be >= 1.
	Threshold int `toml:"threshold" json:"threshold"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Threshold: 3,
	}
}

// Validate returns an error if the Config contains invalid values.
func (c Config) Validate() error {
	if c.Threshold < 1 {
		return fmt.Errorf("anchor: threshold must be >= 1, got %d", c.Threshold)
	}
	return nil
}

// NewFromConfig constructs a Store from a validated Config.
func NewFromConfig(c Config) (*Store, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return New(c.Threshold), nil
}
