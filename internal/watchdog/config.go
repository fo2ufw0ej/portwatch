package watchdog

import (
	"errors"
	"time"
)

// Config holds watchdog configuration.
type Config struct {
	// Timeout is how long without a heartbeat before the daemon is
	// considered stalled. Defaults to 3× the scan interval.
	Timeout time.Duration `toml:"timeout" json:"timeout"`
}

// DefaultConfig returns a sensible default watchdog config.
func DefaultConfig() Config {
	return Config{
		Timeout: 90 * time.Second,
	}
}

// Validate returns an error if the config is invalid.
func (c Config) Validate() error {
	if c.Timeout <= 0 {
		return errors.New("watchdog: timeout must be positive")
	}
	return nil
}

// NewFromConfig constructs a Watchdog from a Config after validation.
func NewFromConfig(cfg Config) (*Watchdog, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return New(cfg.Timeout), nil
}
