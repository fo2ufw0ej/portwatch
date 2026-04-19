package retrier

import (
	"context"
	"errors"
	"time"
)

// Retrier executes a function with retry logic using configurable backoff.
type Retrier struct {
	cfg Config
}

// New returns a Retrier with the given config.
func New(cfg Config) (*Retrier, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Retrier{cfg: cfg}, nil
}

// Do calls fn up to MaxAttempts times, backing off between attempts.
// It stops early if ctx is cancelled or fn returns a non-retryable error.
func (r *Retrier) Do(ctx context.Context, fn func() error) error {
	var last error
	delay := r.cfg.InitialInterval

	for attempt := 1; attempt <= r.cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		last = fn()
		if last == nil {
			return nil
		}

		var nonRetryable *NonRetryableError
		if errors.As(last, &nonRetryable) {
			return last
		}

		if attempt == r.cfg.MaxAttempts {
			break
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}

		delay = time.Duration(float64(delay) * r.cfg.Multiplier)
		if delay > r.cfg.MaxInterval {
			delay = r.cfg.MaxInterval
		}
	}

	return last
}

// NonRetryableError wraps an error that should not be retried.
type NonRetryableError struct {
	Cause error
}

func (e *NonRetryableError) Error() string { return e.Cause.Error() }
func (e *NonRetryableError) Unwrap() error { return e.Cause }

// Permanent marks an error as non-retryable.
func Permanent(err error) error {
	return &NonRetryableError{Cause: err}
}
