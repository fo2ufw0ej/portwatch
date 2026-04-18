package dedup

import "fmt"

// Config holds options for the dedup store.
type Config struct {
	// MaxKeys caps the number of tracked keys to prevent unbounded growth.
	// 0 means unlimited.
	MaxKeys int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{MaxKeys: 256}
}

// Validate returns an error when the config is invalid.
func (c Config) Validate() error {
	if c.MaxKeys < 0 {
		return fmt.Errorf("dedup: MaxKeys must be >= 0, got %d", c.MaxKeys)
	}
	return nil
}

// NewFromConfig constructs a Store, validating the config first.
func NewFromConfig(c Config) (*Store, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return New(), nil
}
