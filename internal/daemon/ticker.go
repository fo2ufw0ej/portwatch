package daemon

import (
	"context"
	"time"
)

// ticker wraps time.Ticker to provide a testable interface for periodic
// execution within the daemon loop.
type ticker interface {
	C() <-chan time.Time
	Stop()
}

// realTicker adapts time.Ticker to satisfy the ticker interface.
type realTicker struct {
	t *time.Ticker
}

// newRealTicker creates a new realTicker that fires at the given interval.
func newRealTicker(d time.Duration) ticker {
	return &realTicker{t: time.NewTicker(d)}
}

// C returns the underlying ticker channel.
func (r *realTicker) C() <-chan time.Time {
	return r.t.C
}

// Stop stops the underlying ticker, releasing associated resources.
func (r *realTicker) Stop() {
	r.t.Stop()
}

// runLoop executes fn on every tick until ctx is cancelled.
// It calls fn once immediately before waiting for the first tick so that
// the daemon produces output without waiting for the first interval to elapse.
func runLoop(ctx context.Context, tk ticker, fn func()){
	// Execute immediately on start so users see output right away.
	fn()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tk.C():
			fn()
		}
	}
}

// mockTicker is a ticker implementation intended for use in tests.
// Callers can manually send on Ch to simulate ticks.
type mockTicker struct {
	Ch chan time.Time
}

// newMockTicker returns a mockTicker with a buffered channel.
func newMockTicker() *mockTicker {
	return &mockTicker{Ch: make(chan time.Time, 1)}
}

// C returns the mock tick channel.
func (m *mockTicker) C() <-chan time.Time {
	return m.Ch
}

// Stop is a no-op for the mock ticker.
func (m *mockTicker) Stop() {}

// Tick sends a tick signal on the mock channel.
func (m *mockTicker) Tick() {
	m.Ch <- time.Now()
}
