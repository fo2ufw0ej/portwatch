// Package grace provides a graceful shutdown coordinator that waits for
// registered components to finish before allowing the process to exit.
package grace

import (
	"context"
	"sync"
	"time"
)

// Coordinator manages graceful shutdown across multiple components.
type Coordinator struct {
	mu      sync.Mutex
	wg      sync.WaitGroup
	cancel  context.CancelFunc
	ctx     context.Context
	timeout time.Duration
}

// New creates a new Coordinator with the given shutdown timeout.
func New(cfg Config) (*Coordinator, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &Coordinator{
		ctx:     ctx,
		cancel:  cancel,
		timeout: cfg.Timeout,
	}, nil
}

// Context returns the coordinator's context, cancelled on shutdown.
func (c *Coordinator) Context() context.Context {
	return c.ctx
}

// Register increments the internal wait group, indicating a component is active.
func (c *Coordinator) Register() {
	c.wg.Add(1)
}

// Done decrements the internal wait group, indicating a component has finished.
func (c *Coordinator) Done() {
	c.wg.Done()
}

// Shutdown signals all components to stop and waits up to the configured
// timeout for them to finish. Returns true if all components finished cleanly.
func (c *Coordinator) Shutdown() bool {
	c.cancel()

	finished := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
		return true
	case <-time.After(c.timeout):
		return false
	}
}
