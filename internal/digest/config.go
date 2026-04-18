package digest

import (
	"errors"
	"os"
	"path/filepath"
)

// Config holds configuration for the digest store.
type Config struct {
	// StorePath is the file path where the digest is persisted.
	StorePath string `toml:"store_path"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}
	return Config{
		StorePath: filepath.Join(cacheDir, "portwatch", "digest.json"),
	}
}

// Validate returns an error if the config is invalid.
func (c Config) Validate() error {
	if c.StorePath == "" {
		return errors.New("digest: store_path must not be empty")
	}
	return nil
}

// NewStore builds a Store from the config.
func NewStoreFromConfig(c Config) (*Store, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return NewStore(c.StorePath), nil
}
