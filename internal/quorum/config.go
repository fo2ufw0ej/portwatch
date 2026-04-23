package quorum

import "fmt"

// NewFromConfig constructs a Quorum from a Config, returning an error
// if validation fails.
func NewFromConfig(cfg Config) (*Quorum, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("quorum: invalid config: %w", err)
	}
	return New(cfg)
}
