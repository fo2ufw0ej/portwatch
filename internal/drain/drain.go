// Package drain provides a graceful event-draining mechanism that flushes
// buffered items before shutdown, ensuring no events are silently dropped.
package drain

import (
	"context"
	"sync"
	"time"
)

// Config holds tunable parameters for the Drainer.
type Config struct {
	// Capacity is the maximum number of items buffered before blocking.
	Capacity int
	// FlushTimeout is the maximum time to wait for all items to be consumed
	// during a Drain call.
	FlushTimeout time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Capacity:     256,
		FlushTimeout: 5 * time.Second,
	}
}

// Validate returns an error if the Config is invalid.
func (c Config) Validate() error {
	if c.Capacity <= 0 {
		return ErrInvalidCapacity
	}
	if c.FlushTimeout <= 0 {
		return ErrInvalidFlushTimeout
	}
	return nil
}

// Drainer buffers items and allows a consumer goroutine to process them.
// Call Drain to block until all buffered items are consumed or the flush
// timeout elapses.
type Drainer[T any] struct {
	cfg Config
	ch  chan T
	wg  sync.WaitGroup
}

// New creates a new Drainer with the given Config.
func New[T any](cfg Config) (*Drainer[T], error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Drainer[T]{
		cfg: cfg,
		ch:  make(chan T, cfg.Capacity),
	}, nil
}

// Push enqueues an item. It returns false if the context is done before the
// item could be enqueued.
func (d *Drainer[T]) Push(ctx context.Context, item T) bool {
	select {
	case d.ch <- item:
		d.wg.Add(1)
		return true
	case <-ctx.Done():
		return false
	}
}

// C returns the read channel that a consumer should range over.
func (d *Drainer[T]) C() <-chan T {
	return d.ch
}

// Done signals that one item has been fully processed. Must be called once
// per successful Push.
func (d *Drainer[T]) Done() {
	d.wg.Done()
}

// Drain closes the input channel and waits for all pushed items to be marked
// done, or until the flush timeout expires. It returns true if all items were
// consumed in time.
func (d *Drainer[T]) Drain() bool {
	close(d.ch)
	finished := make(chan struct{})
	go func() {
		d.wg.Wait()
		close(finished)
	}()
	select {
	case <-finished:
		return true
	case <-time.After(d.cfg.FlushTimeout):
		return false
	}
}
