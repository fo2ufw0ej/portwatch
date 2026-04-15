package alert

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/ratelimit"
	"github.com/user/portwatch/internal/scanner"
)

// RateLimitedNotifier wraps a Notifier and suppresses repeated alerts
// using a token-bucket rate limiter keyed on change direction.
type RateLimitedNotifier struct {
	inner   *Notifier
	limiter *ratelimit.Limiter
}

// NewRateLimitedNotifier creates a RateLimitedNotifier that allows at most
// maxBurst alert events per direction ("opened" / "closed") within window.
func NewRateLimitedNotifier(inner *Notifier, window time.Duration, maxBurst int) *RateLimitedNotifier {
	return &RateLimitedNotifier{
		inner:   inner,
		limiter: ratelimit.New(window, maxBurst),
	}
}

// Notify forwards diff to the inner Notifier only when the rate limiter
// permits. Suppressed alerts are counted but not forwarded.
func (r *RateLimitedNotifier) Notify(diff scanner.Diff) error {
	var errs []error

	if len(diff.Opened) > 0 {
		if r.limiter.Allow("opened") {
			if err := r.inner.Notify(scanner.Diff{Opened: diff.Opened}); err != nil {
				errs = append(errs, fmt.Errorf("opened notify: %w", err))
			}
		}
	}

	if len(diff.Closed) > 0 {
		if r.limiter.Allow("closed") {
			if err := r.inner.Notify(scanner.Diff{Closed: diff.Closed}); err != nil {
				errs = append(errs, fmt.Errorf("closed notify: %w", err))
			}
		}
	}

	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
