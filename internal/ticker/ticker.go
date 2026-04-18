// Package ticker provides a configurable interval ticker with jitter support.
package ticker

import (
	"fmt"
	"math/rand"
	"time"
)

// Config holds ticker configuration.
type Config struct {
	Interval time.Duration
	Jitter   float64 // fraction of Interval to add as random jitter [0, 1]
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Interval: 30 * time.Second,
		Jitter:   0.1,
	}
}

// Validate returns an error if the config is invalid.
func (c Config) Validate() error {
	if c.Interval <= 0 {
		return fmt.Errorf("ticker: interval must be positive, got %v", c.Interval)
	}
	if c.Jitter < 0 || c.Jitter > 1 {
		return fmt.Errorf("ticker: jitter must be in [0, 1], got %v", c.Jitter)
	}
	return nil
}

// Ticker fires at a configurable interval with optional jitter.
type Ticker struct {
	cfg  Config
	ch   chan time.Time
	stop chan struct{}
}

// New creates and starts a new Ticker.
func New(cfg Config) (*Ticker, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	t := &Ticker{
		cfg:  cfg,
		ch:   make(chan time.Time, 1),
		stop: make(chan struct{}),
	}
	go t.run()
	return t, nil
}

// C returns the channel on which ticks are delivered.
func (t *Ticker) C() <-chan time.Time {
	return t.ch
}

// Stop halts the ticker.
func (t *Ticker) Stop() {
	close(t.stop)
}

func (t *Ticker) run() {
	for {
		delay := t.cfg.Interval
		if t.cfg.Jitter > 0 {
			jitter := time.Duration(float64(t.cfg.Interval) * t.cfg.Jitter * rand.Float64())
			delay += jitter
		}
		select {
		case <-time.After(delay):
			select {
			case t.ch <- time.Now():
			default:
			}
		case <-t.stop:
			return
		}
	}
}
