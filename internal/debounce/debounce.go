// Package debounce provides a simple debouncer that delays action execution
// until a quiet period has elapsed, preventing alert storms on flapping ports.
package debounce

import (
	"sync"
	"time"
)

// Debouncer holds pending triggers per key and fires them after a quiet window.
type Debouncer struct {
	mu      sync.Mutex
	delay   time.Duration
	timers  map[string]*time.Timer
	callback func(key string)
}

// New creates a Debouncer that waits delay before invoking callback for a key.
func New(delay time.Duration, callback func(key string)) *Debouncer {
	return &Debouncer{
		delay:    delay,
		timers:   make(map[string]*time.Timer),
		callback: callback,
	}
}

// Trigger schedules callback(key) after the debounce delay.
// If called again before the delay elapses, the timer resets.
func (d *Debouncer) Trigger(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[key]; ok {
		t.Reset(d.delay)
		return
	}

	d.timers[key] = time.AfterFunc(d.delay, func() {
		d.mu.Lock()
		delete(d.timers, key)
		d.mu.Unlock()
		d.callback(key)
	})
}

// Cancel cancels a pending trigger for key, if any.
func (d *Debouncer) Cancel(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[key]; ok {
		t.Stop()
		delete(d.timers, key)
	}
}

// Pending returns the number of keys with active pending timers.
func (d *Debouncer) Pending() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.timers)
}
