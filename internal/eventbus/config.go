package eventbus

import "fmt"

// Config holds optional settings for the event bus.
type Config struct {
	// BufferSize is reserved for future async dispatch support.
	// A value of 0 means synchronous (current behaviour).
	BufferSize int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{BufferSize: 0}
}

// Validate returns an error if the Config is invalid.
func (c Config) Validate() error {
	if c.BufferSize < 0 {
		return fmt.Errorf("eventbus: BufferSize must be >= 0, got %d", c.BufferSize)
	}
	return nil
}
