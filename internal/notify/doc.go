// Package notify provides webhook-based alerting for portwatch.
//
// When a port change is detected by the daemon, a WebhookNotifier can POST
// a JSON payload to a configured HTTP endpoint. The payload includes the
// lists of newly opened and closed ports along with a timestamp and a
// human-readable message.
//
// Usage:
//
//	cfg := notify.DefaultConfig()
//	cfg.Enabled = true
//	cfg.Webhook.URL = "https://hooks.example.com/portwatch"
//	if err := cfg.Validate(); err != nil { ... }
//	wn, err := notify.NewWebhookNotifier(cfg.Webhook)
//	if err != nil { ... }
//	wn.Send(notify.WebhookPayload{ ... })
package notify
