package scorecard

import "fmt"

// Config controls scorecard behaviour.
type Config struct {
	// MinScore is the threshold below which a port is considered unhealthy.
	MinScore float64
	// MaxPorts caps how many ports the card will track (0 = unlimited).
	MaxPorts int
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		MinScore: 0.5,
		MaxPorts: 0,
	}
}

// Validate returns an error if the config is invalid.
func (c Config) Validate() error {
	if c.MinScore < 0 || c.MinScore > 1 {
		return fmt.Errorf("scorecard: MinScore must be in [0, 1], got %.2f", c.MinScore)
	}
	if c.MaxPorts < 0 {
		return fmt.Errorf("scorecard: MaxPorts must be >= 0, got %d", c.MaxPorts)
	}
	return nil
}

// Unhealthy returns true when the given score falls below the configured threshold.
func (c Config) Unhealthy(score float64) bool {
	return score < c.MinScore
}
