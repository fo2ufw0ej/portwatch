// Package circuitbreaker provides a simple circuit breaker to suppress
// repeated alert delivery attempts when a downstream notifier is failing.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// State represents the circuit breaker state.
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

// ErrOpen is returned when the circuit is open.
var ErrOpen = errors.New("circuit breaker is open")

// Breaker is a simple circuit breaker.
type Breaker struct {
	mu          sync.Mutex
	maxFailures int
	resetAfter  time.Duration
	failures    int
	lastFailure time.Time
	state       State
}

// New creates a new Breaker. After maxFailures consecutive failures the circuit
// opens and no calls are allowed until resetAfter has elapsed.
func New(maxFailures int, resetAfter time.Duration) *Breaker {
	return &Breaker{
		maxFailures: maxFailures,
		resetAfter:  resetAfter,
	}
}

// Allow returns nil if the call should proceed, or ErrOpen if the circuit is open.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.state == StateOpen {
		if time.Since(b.lastFailure) >= b.resetAfter {
			b.state = StateHalfOpen
			return nil
		}
		return ErrOpen
	}
	return nil
}

// RecordSuccess resets the breaker to closed.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure records a failure and opens the circuit if the threshold is met.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	b.lastFailure = time.Now()
	if b.failures >= b.maxFailures {
		b.state = StateOpen
	}
}

// State returns the current state.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
