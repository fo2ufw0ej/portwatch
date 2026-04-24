package ledger

import (
	"errors"
	"fmt"
)

// Config holds configuration for the Ledger.
type Config struct {
	// Path is the file path where ledger data is persisted.
	Path string `toml:"path" json:"path"`

	// FlappingThreshold is the minimum number of state changes (open+close)
	// before a port is considered flapping.
	FlappingThreshold int `toml:"flapping_threshold" json:"flapping_threshold"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Path:              "/var/lib/portwatch/ledger.json",
		FlappingThreshold: 5,
	}
}

// Validate returns an error if the Config is invalid.
func (c Config) Validate() error {
	if c.Path == "" {
		return errors.New("ledger: path must not be empty")
	}
	if c.FlappingThreshold < 1 {
		return fmt.Errorf("ledger: flapping_threshold must be >= 1, got %d", c.FlappingThreshold)
	}
	return nil
}

// NewFromConfig creates a Ledger from a validated Config.
func NewFromConfig(c Config) (*Ledger, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return New(c.Path)
}
