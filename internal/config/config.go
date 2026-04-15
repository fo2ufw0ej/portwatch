// Package config loads and validates portwatch configuration.
package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all runtime configuration for portwatch.
type Config struct {
	Interval    time.Duration `yaml:"interval"`
	PortRange   PortRange     `yaml:"port_range"`
	AlertWindow time.Duration `yaml:"alert_window"`
	AlertBurst  int           `yaml:"alert_burst"`
	StatePath   string        `yaml:"state_path"`
	LogPath     string        `yaml:"log_path"`
	Format      string        `yaml:"format"`
}

// PortRange defines the inclusive [Low, High] port scan range.
type PortRange struct {
	Low  int `yaml:"low"`
	High int `yaml:"high"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Interval:    30 * time.Second,
		PortRange:   PortRange{Low: 1, High: 1024},
		AlertWindow: time.Minute,
		AlertBurst:  5,
		StatePath:   "/tmp/portwatch_state.json",
		LogPath:     "/tmp/portwatch_history.json",
		Format:      "text",
	}
}

// Load reads a YAML config file from path and merges it over defaults.
// Returns an error if the file cannot be read or values are invalid.
func Load(path string) (Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("read config: %w", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parse config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return cfg, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

func (c Config) validate() error {
	if c.Interval <= 0 {
		return errors.New("interval must be positive")
	}
	if c.PortRange.Low < 1 || c.PortRange.High > 65535 || c.PortRange.Low > c.PortRange.High {
		return fmt.Errorf("invalid port range [%d, %d]", c.PortRange.Low, c.PortRange.High)
	}
	if c.AlertBurst < 1 {
		return errors.New("alert_burst must be at least 1")
	}
	if c.AlertWindow <= 0 {
		return errors.New("alert_window must be positive")
	}
	if c.Format != "text" && c.Format != "json" {
		return fmt.Errorf("unsupported format %q: must be text or json", c.Format)
	}
	return nil
}
