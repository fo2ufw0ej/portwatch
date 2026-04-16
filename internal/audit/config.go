package audit

import (
	"fmt"
	"io"
	"os"
)

// Config holds configuration for the audit logger.
type Config struct {
	Enabled bool   `toml:"enabled"`
	Path    string `toml:"path"`
	Format  string `toml:"format"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Enabled: false,
		Path:    "",
		Format:  "text",
	}
}

// Validate returns an error if the config is invalid.
func (c Config) Validate() error {
	if c.Format != "text" && c.Format != "json" {
		return fmt.Errorf("audit: invalid format %q, must be \"text\" or \"json\"", c.Format)
	}
	return nil
}

// Open returns an io.WriteCloser for the configured audit path.
// If Path is empty, os.Stdout is returned with a no-op closer.
func (c Config) Open() (io.WriteCloser, error) {
	if c.Path == "" {
		return nopCloser{os.Stdout}, nil
	}
	f, err := os.OpenFile(c.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	return f, nil
}

type nopCloser struct{ io.Writer }

func (nopCloser) Close() error { return nil }
