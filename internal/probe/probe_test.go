package probe_test

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/user/portwatch/internal/probe"
)

func startListener(t *testing.T) (port int, stop func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("startListener: %v", err)
	}
	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	p, _ := strconv.Atoi(portStr)
	return p, func() { _ = ln.Close() }
}

func defaultProbe(t *testing.T) *probe.Probe {
	t.Helper()
	p, err := probe.New(probe.DefaultConfig())
	if err != nil {
		t.Fatalf("probe.New: %v", err)
	}
	return p
}

func TestCheck_OpenPort_Reachable(t *testing.T) {
	port, stop := startListener(t)
	defer stop()

	p := defaultProbe(t)
	res := p.Check(context.Background(), port)

	if !res.Reachable {
		t.Fatalf("expected reachable, got err: %v", res.Err)
	}
	if res.Latency <= 0 {
		t.Error("expected positive latency")
	}
	if p.Failures(port) != 0 {
		t.Errorf("expected 0 failures, got %d", p.Failures(port))
	}
}

func TestCheck_ClosedPort_NotReachable(t *testing.T) {
	p := defaultProbe(t)
	res := p.Check(context.Background(), 1) // port 1 is never open in tests

	if res.Reachable {
		t.Fatal("expected unreachable")
	}
	if res.Err == nil {
		t.Error("expected non-nil error")
	}
	if p.Failures(1) != 1 {
		t.Errorf("expected 1 failure, got %d", p.Failures(1))
	}
}

func TestCheck_AccumulatesFailures(t *testing.T) {
	p := defaultProbe(t)
	for i := 1; i <= 3; i++ {
		p.Check(context.Background(), 1)
		if p.Failures(1) != i {
			t.Errorf("iteration %d: expected %d failures, got %d", i, i, p.Failures(1))
		}
	}
}

func TestReset_ClearsFailures(t *testing.T) {
	p := defaultProbe(t)
	p.Check(context.Background(), 1)
	p.Check(context.Background(), 1)
	p.Reset(1)
	if p.Failures(1) != 0 {
		t.Errorf("expected 0 after reset, got %d", p.Failures(1))
	}
}

func TestNew_InvalidConfig(t *testing.T) {
	cfg := probe.DefaultConfig()
	cfg.Timeout = -1 * time.Second
	_, err := probe.New(cfg)
	if err == nil {
		t.Fatal("expected error for negative timeout")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := probe.DefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("DefaultConfig should be valid: %v", err)
	}
}
