// Package watchdog detects and reports stale scan cycles.
package watchdog

import (
	"sync"
	"time"
)

// Watchdog tracks the last successful scan time and exposes
// a method to check whether the daemon appears stalled.
type Watchdog struct {
	mu       sync.Mutex
	lastBeat time.Time
	timeout  time.Duration
}

// New creates a Watchdog with the given stall timeout.
func New(timeout time.Duration) *Watchdog {
	return &Watchdog{
		timeout:  timeout,
		lastBeat: time.Now(),
	}
}

// Beat records a successful scan heartbeat.
func (w *Watchdog) Beat() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.lastBeat = time.Now()
}

// Stalled returns true if no Beat has been recorded within the timeout window.
func (w *Watchdog) Stalled() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return time.Since(w.lastBeat) > w.timeout
}

// LastBeat returns the time of the most recent heartbeat.
func (w *Watchdog) LastBeat() time.Time {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.lastBeat
}

// StalledFor returns how long the watchdog has been stalled, or zero.
func (w *Watchdog) StalledFor() time.Duration {
	w.mu.Lock()
	defer w.mu.Unlock()
	elapsed := time.Since(w.lastBeat)
	if elapsed <= w.timeout {
		return 0
	}
	return elapsed - w.timeout
}
