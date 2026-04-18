package cooldown

import (
	"errors"
	"time"
)

// Config holds configuration for a Cooldown instance.
type Config struct {
	Period time.Duration `json:"period"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Period: 30 * time.Second,
	}
}

// Validate returns an error if the config is invalid.
func (c Config) Validate() error {
	if c.Period <= 0 {
		return errors.New("cooldown: period must be positive")
	}
	return nil
}

// NewFromConfig constructs a Cooldown from a Config.
func NewFromConfig(cfg Config) (*Cooldown, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return New(cfg.Period), nil
}
