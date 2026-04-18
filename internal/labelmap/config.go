package labelmap

import "fmt"

// Config holds configuration for the LabelMap.
type Config struct {
	// Custom maps port numbers to user-defined labels.
	Custom map[int]string `json:"custom"`
}

// DefaultConfig returns a Config with no custom labels.
func DefaultConfig() Config {
	return Config{
		Custom: map[int]string{},
	}
}

// Validate checks that the Config is valid.
func (c Config) Validate() error {
	for port := range c.Custom {
		if port < 1 || port > 65535 {
			return fmt.Errorf("labelmap: invalid port %d in custom labels", port)
		}
	}
	return nil
}

// NewFromConfig constructs a LabelMap from Config.
func NewFromConfig(c Config) (*LabelMap, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return New(c.Custom), nil
}
