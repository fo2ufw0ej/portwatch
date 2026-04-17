package signal

// Config holds configuration for signal handling behaviour.
type Config struct {
	// Signals lists the OS signals that trigger shutdown.
	// Defaults to SIGINT and SIGTERM when empty.
	Signals []string `toml:"signals" json:"signals"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Signals: []string{"SIGINT", "SIGTERM"},
	}
}

// Validate checks that Config values are acceptable.
func (c Config) Validate() error {
	if len(c.Signals) == 0 {
		return errNoSignals
	}
	return nil
}

import "errors"

var errNoSignals = errors.New("signal: at least one signal must be specified")
