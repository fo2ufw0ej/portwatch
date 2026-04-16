package notify

import (
	"fmt"
	"time"
)

// Config holds top-level notification configuration.
type Config struct {
	Webhook WebhookConfig
	Enabled bool
}

// DefaultConfig returns a Config with notifications disabled.
func DefaultConfig() Config {
	return Config{
		Enabled: false,
		Webhook: WebhookConfig{
			Timeout: 5 * time.Second,
		},
	}
}

// Validate checks that the config is internally consistent.
func (c Config) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Webhook.URL == "" {
		return fmt.Errorf("notify: webhook URL required when notifications are enabled")
	}
	if c.Webhook.Timeout <= 0 {
		return fmt.Errorf("notify: webhook timeout must be positive")
	}
	return nil
}
