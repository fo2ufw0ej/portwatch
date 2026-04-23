package probe

import (
	"errors"
	"time"
)

// Config holds tunable parameters for Probe.
type Config struct {
	// Host is the target hostname or IP address.
	Host string
	// Timeout is the maximum time to wait for a connection.
	Timeout time.Duration
	// MaxFailures is the threshold at which a port is considered degraded.
	MaxFailures int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Host:        "127.0.0.1",
		Timeout:     2 * time.Second,
		MaxFailures: 3,
	}
}

// Validate returns an error if the Config is invalid.
func (c Config) Validate() error {
	if c.Host == "" {
		return errors.New("host must not be empty")
	}
	if c.Timeout <= 0 {
		return errors.New("timeout must be positive")
	}
	if c.MaxFailures < 1 {
		return errors.New("max_failures must be at least 1")
	}
	return nil
}

// NewFromConfig returns a Probe built from cfg.
func NewFromConfig(cfg Config) (*Probe, error) {
	return New(cfg)
}
