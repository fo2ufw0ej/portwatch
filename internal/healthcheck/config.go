package healthcheck

import (
	"errors"
	"fmt"
	"net"
)

// Config holds configuration for the health check HTTP server.
type Config struct {
	// Enabled controls whether the health server is started.
	Enabled bool `toml:"enabled" json:"enabled"`
	// Host is the interface to bind to (default: 127.0.0.1).
	Host string `toml:"host" json:"host"`
	// Port is the TCP port to listen on (default: 9090).
	Port int `toml:"port" json:"port"`
}

// DefaultConfig returns a Config with safe defaults.
func DefaultConfig() Config {
	return Config{
		Enabled: false,
		Host:    "127.0.0.1",
		Port:    9090,
	}
}

// Validate checks that the Config fields are valid.
func (c Config) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("healthcheck: port %d out of range (1-65535)", c.Port)
	}
	if c.Host == "" {
		return errors.New("healthcheck: host must not be empty")
	}
	if net.ParseIP(c.Host) == nil {
		// Allow hostnames too; only reject clearly empty strings.
		if len(c.Host) == 0 {
			return errors.New("healthcheck: invalid host")
		}
	}
	return nil
}

// Addr returns the combined host:port address string.
func (c Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
