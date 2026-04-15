package healthcheck_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/user/portwatch/internal/healthcheck"
	"github.com/user/portwatch/internal/metrics"
)

func freePort(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("could not find free port: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	_ = l.Close()
	return fmt.Sprintf("127.0.0.1:%d", port)
}

func TestHealthEndpoint_ReturnsOK(t *testing.T) {
	m := metrics.New()
	addr := freePort(t)
	srv := healthcheck.New(addr, m)
	_ = srv.Start()
	defer srv.Stop()
	time.Sleep(30 * time.Millisecond)

	resp, err := http.Get("http://" + addr + "/health")
	if err != nil {
		t.Fatalf("GET /health: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != `{"status":"ok"}` {
		t.Errorf("unexpected body: %s", body)
	}
}

func TestMetricsEndpoint_ReturnsJSON(t *testing.T) {
	m := metrics.New()
	m.RecordScan(3)
	addr := freePort(t)
	srv := healthcheck.New(addr, m)
	_ = srv.Start()
	defer srv.Stop()
	time.Sleep(30 * time.Millisecond)

	resp, err := http.Get("http://" + addr + "/metrics")
	if err != nil {
		t.Fatalf("GET /metrics: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if _, ok := result["total_scans"]; !ok {
		t.Error("expected total_scans field in metrics response")
	}
}
