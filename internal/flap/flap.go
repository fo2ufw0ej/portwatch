// Package flap detects ports that are rapidly toggling between open and
// closed states ("flapping"), which often indicates unstable services or
// misconfigured firewall rules.
package flap

import (
	"sync"
	"time"
)

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Window:    2 * time.Minute,
		Threshold: 4,
	}
}

// Config controls flap detection behaviour.
type Config struct {
	// Window is the rolling duration over which state changes are counted.
	Window time.Duration
	// Threshold is the minimum number of state changes within Window that
	// classifies a port as flapping.
	Threshold int
}

// Validate returns an error if the Config is invalid.
func (c Config) Validate() error {
	if c.Window <= 0 {
		return ErrInvalidWindow
	}
	if c.Threshold < 2 {
		return ErrInvalidThreshold
	}
	return nil
}

// event records a single open/close transition for a port.
type event struct {
	at time.Time
}

// Detector tracks state changes per port and reports flapping.
type Detector struct {
	cfg    Config
	mu     sync.Mutex
	events map[int][]event
}

// New creates a Detector from cfg. It returns an error if cfg is invalid.
func New(cfg Config) (*Detector, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Detector{
		cfg:    cfg,
		events: make(map[int][]event),
	}, nil
}

// Record registers a state change for port at the given time.
func (d *Detector) Record(port int, at time.Time) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.prune(port, at)
	d.events[port] = append(d.events[port], event{at: at})
}

// Flapping reports whether port has exceeded the change threshold within the
// configured window as of now.
func (d *Detector) Flapping(port int, now time.Time) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.prune(port, now)
	return len(d.events[port]) >= d.cfg.Threshold
}

// Reset clears all recorded events for port.
func (d *Detector) Reset(port int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.events, port)
}

// prune removes events older than Window relative to now. Must be called with
// d.mu held.
func (d *Detector) prune(port int, now time.Time) {
	cutoff := now.Add(-d.cfg.Window)
	evs := d.events[port]
	var kept []event
	for _, e := range evs {
		if !e.at.Before(cutoff) {
			kept = append(kept, e)
		}
	}
	if len(kept) == 0 {
		delete(d.events, port)
	} else {
		d.events[port] = kept
	}
}
