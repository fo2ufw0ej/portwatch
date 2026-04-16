package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookConfig holds configuration for a webhook notifier.
type WebhookConfig struct {
	URL     string
	Timeout time.Duration
	Secret  string // optional: added as X-PortWatch-Secret header
}

// DefaultWebhookConfig returns a WebhookConfig with sensible defaults.
func DefaultWebhookConfig() WebhookConfig {
	return WebhookConfig{
		Timeout: 5 * time.Second,
	}
}

// WebhookPayload is the JSON body sent on each alert.
type WebhookPayload struct {
	Timestamp string   `json:"timestamp"`
	Opened    []int    `json:"opened"`
	Closed    []int    `json:"closed"`
	Message   string   `json:"message"`
}

// WebhookNotifier sends alert payloads to an HTTP endpoint.
type WebhookNotifier struct {
	cfg    WebhookConfig
	client *http.Client
}

// NewWebhookNotifier creates a WebhookNotifier from cfg.
func NewWebhookNotifier(cfg WebhookConfig) (*WebhookNotifier, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("webhook URL must not be empty")
	}
	return &WebhookNotifier{
		cfg:    cfg,
		client: &http.Client{Timeout: cfg.Timeout},
	}, nil
}

// Send posts payload to the configured webhook URL.
func (w *WebhookNotifier) Send(payload WebhookPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, w.cfg.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if w.cfg.Secret != "" {
		req.Header.Set("X-PortWatch-Secret", w.cfg.Secret)
	}
	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("send webhook: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}
	return nil
}
