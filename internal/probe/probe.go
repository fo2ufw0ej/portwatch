// Package probe provides a periodic connectivity probe that verifies
// individual ports remain reachable and tracks consecutive failure counts.
package probe

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

// Result holds the outcome of a single probe attempt.
type Result struct {
	Port    int
	Reachable bool
	Latency time.Duration
	Err     error
}

// Probe checks whether a TCP port is reachable within the given timeout.
type Probe struct {
	cfg     Config
	mu      sync.Mutex
	failures map[int]int
}

// New creates a Probe with the provided configuration.
func New(cfg Config) (*Probe, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("probe: invalid config: %w", err)
	}
	return &Probe{
		cfg:      cfg,
		failures: make(map[int]int),
	}, nil
}

// Check performs a TCP dial against host:port and returns a Result.
func (p *Probe) Check(ctx context.Context, port int) Result {
	addr := fmt.Sprintf("%s:%d", p.cfg.Host, port)
	start := time.Now()

	conn, err := (&net.Dialer{Timeout: p.cfg.Timeout}).DialContext(ctx, "tcp", addr)
	latency := time.Since(start)

	p.mu.Lock()
	defer p.mu.Unlock()

	if err != nil {
		p.failures[port]++
		return Result{Port: port, Reachable: false, Latency: latency, Err: err}
	}
	_ = conn.Close()
	p.failures[port] = 0
	return Result{Port: port, Reachable: true, Latency: latency}
}

// Failures returns the consecutive failure count for a port.
func (p *Probe) Failures(port int) int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.failures[port]
}

// Reset clears the failure counter for a port.
func (p *Probe) Reset(port int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.failures, port)
}
