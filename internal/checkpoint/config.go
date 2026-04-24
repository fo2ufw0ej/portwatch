package checkpoint

import (
	"fmt"
	"time"
)

// Config controls checkpoint behaviour.
type Config struct {
	// Path is the file path where checkpoints are persisted.
	Path string `json:"path"`

	// Interval is how often the daemon writes a new checkpoint.
	Interval time.Duration `json:"interval"`

	// Enabled toggles checkpointing entirely.
	Enabled bool `json:"enabled"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Path:     "/var/lib/portwatch/checkpoint.json",
		Interval: 5 * time.Minute,
		Enabled:  true,
	}
}

// Validate returns an error if the configuration is invalid.
func (c Config) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Path == "" {
		return fmt.Errorf("checkpoint: path must not be empty")
	}
	if c.Interval <= 0 {
		return fmt.Errorf("checkpoint: interval must be positive, got %v", c.Interval)
	}
	return nil
}

// NewFromConfig returns a Store using the supplied Config.
// If cfg.Enabled is false, NewFromConfig returns nil, nil.
func NewFromConfig(cfg Config) (*Store, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if !cfg.Enabled {
		return nil, nil
	}
	return New(cfg.Path)
}
