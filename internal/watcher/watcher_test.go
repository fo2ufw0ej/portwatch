package watcher_test

import (
	"net"
	"testing"
	"time"

	"github.com/example/portwatch/internal/watcher"
)

// startListener opens a TCP listener on a random port and returns it.
func startListener(t *testing.T) (net.Listener, int) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return ln, port
}

func TestProbe_OpenPort(t *testing.T) {
	ln, port := startListener(t)
	defer ln.Close()

	status := watcher.Probe("127.0.0.1", port, time.Second)

	if !status.Open {
		t.Errorf("expected port %d to be open", port)
	}
	if status.Port != port {
		t.Errorf("expected Port=%d, got %d", port, status.Port)
	}
	if status.Latency <= 0 {
		t.Errorf("expected positive latency, got %v", status.Latency)
	}
	if status.CheckedAt.IsZero() {
		t.Error("expected CheckedAt to be set")
	}
}

func TestProbe_ClosedPort(t *testing.T) {
	// Use a port that is almost certainly not listening.
	status := watcher.Probe("127.0.0.1", 1, 200*time.Millisecond)

	if status.Open {
		t.Error("expected port 1 to be closed")
	}
}

func TestProbeAll_MixedPorts(t *testing.T) {
	ln, openPort := startListener(t)
	defer ln.Close()

	ports := []int{openPort, 1}
	results := watcher.ProbeAll("127.0.0.1", ports, 500*time.Millisecond)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if !results[openPort].Open {
		t.Errorf("expected port %d to be open", openPort)
	}
	if results[1].Open {
		t.Error("expected port 1 to be closed")
	}
}
