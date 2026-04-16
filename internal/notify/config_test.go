package notify_test

import (
	"testing"
	"time"

	"github.com/rednexie/portwatch/internal/notify"
)

func TestDefaultConfig_Disabled(t *testing.T) {
	cfg := notify.DefaultConfig()
	if cfg.Enabled {
		t.Error("expected notifications disabled by default")
	}
	if cfg.Webhook.Timeout != 5*time.Second {
		t.Errorf("unexpected default timeout: %v", cfg.Webhook.Timeout)
	}
}

func TestValidate_DisabledSkipsURLCheck(t *testing.T) {
	cfg := notify.DefaultConfig()
	cfg.Enabled = false
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestValidate_EnabledNoURL(t *testing.T) {
	cfg := notify.DefaultConfig()
	cfg.Enabled = true
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing URL")
	}
}

func TestValidate_EnabledWithURL(t *testing.T) {
	cfg := notify.DefaultConfig()
	cfg.Enabled = true
	cfg.Webhook.URL = "http://example.com/hook"
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_ZeroTimeout(t *testing.T) {
	cfg := notify.DefaultConfig()
	cfg.Enabled = true
	cfg.Webhook.URL = "http://example.com/hook"
	cfg.Webhook.Timeout = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero timeout")
	}
}
