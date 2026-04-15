package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the portwatch daemon configuration.
type Config struct {
	ScanInterval time.Duration `yaml:"-"`
	ScanIntervalRaw string     `yaml:"scan_interval"`
	PortRange       PortRange  `yaml:"port_range"`
	AlertOutput     string     `yaml:"alert_output"`
	LogLevel        string     `yaml:"log_level"`
}

// PortRange defines the inclusive range of ports to scan.
type PortRange struct {
	From int `yaml:"from"`
	To   int `yaml:"to"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		ScanInterval:    30 * time.Second,
		ScanIntervalRaw: "30s",
		PortRange: PortRange{
			From: 1,
			To:   65535,
		},
		AlertOutput: "stdout",
		LogLevel:    "info",
	}
}

// Load reads a YAML config file from path and returns a validated Config.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if cfg.ScanIntervalRaw != "" {
		d, err := time.ParseDuration(cfg.ScanIntervalRaw)
		if err != nil {
			return nil, fmt.Errorf("invalid scan_interval %q: %w", cfg.ScanIntervalRaw, err)
		}
		cfg.ScanInterval = d
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.PortRange.From < 1 || c.PortRange.From > 65535 {
		return fmt.Errorf("port_range.from must be between 1 and 65535, got %d", c.PortRange.From)
	}
	if c.PortRange.To < 1 || c.PortRange.To > 65535 {
		return fmt.Errorf("port_range.to must be between 1 and 65535, got %d", c.PortRange.To)
	}
	if c.PortRange.From > c.PortRange.To {
		return fmt.Errorf("port_range.from (%d) must be <= port_range.to (%d)", c.PortRange.From, c.PortRange.To)
	}
	if c.ScanInterval <= 0 {
		return fmt.Errorf("scan_interval must be positive")
	}
	return nil
}
