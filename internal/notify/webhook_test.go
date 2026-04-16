package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rednexie/portwatch/internal/notify"
)

func TestNewWebhookNotifier_EmptyURL(t *testing.T) {
	_, err := notify.NewWebhookNotifier(notify.WebhookConfig{})
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestSend_Success(t *testing.T) {
	var received notify.WebhookPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	wn, err := notify.NewWebhookNotifier(notify.WebhookConfig{URL: ts.URL, Timeout: 2 * time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	payload := notify.WebhookPayload{
		Timestamp: "2024-01-01T00:00:00Z",
		Opened:    []int{8080},
		Closed:    []int{},
		Message:   "port opened",
	}
	if err := wn.Send(payload); err != nil {
		t.Fatalf("Send failed: %v", err)
	}
	if len(received.Opened) != 1 || received.Opened[0] != 8080 {
		t.Errorf("unexpected payload: %+v", received)
	}
}

func TestSend_SecretHeader(t *testing.T) {
	var gotSecret string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotSecret = r.Header.Get("X-PortWatch-Secret")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	wn, _ := notify.NewWebhookNotifier(notify.WebhookConfig{URL: ts.URL, Secret: "mysecret", Timeout: 2 * time.Second})
	wn.Send(notify.WebhookPayload{})
	if gotSecret != "mysecret" {
		t.Errorf("expected secret header, got %q", gotSecret)
	}
}

func TestSend_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	wn, _ := notify.NewWebhookNotifier(notify.WebhookConfig{URL: ts.URL, Timeout: 2 * time.Second})
	if err := wn.Send(notify.WebhookPayload{}); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}
