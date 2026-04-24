package roster

import (
	"errors"
	"time"
)

// Config controls optional behaviour of a managed Roster.
type Config struct {
	// MaxAge is the duration after which an entry that has not been
	// touched is considered stale and eligible for eviction.
	// Zero disables age-based eviction.
	MaxAge time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxAge: 5 * time.Minute,
	}
}

// Validate returns an error if the Config is invalid.
func (c Config) Validate() error {
	if c.MaxAge < 0 {
		return errors.New("roster: MaxAge must be non-negative")
	}
	return nil
}

// Evict removes all entries whose LastSeen is older than MaxAge.
// It is a no-op when MaxAge is zero.
func (r *Roster) Evict(cfg Config) int {
	if cfg.MaxAge == 0 {
		return 0
	}
	cutoff := r.now().Add(-cfg.MaxAge)
	r.mu.Lock()
	defer r.mu.Unlock()
	removed := 0
	for p, e := range r.entries {
		if e.LastSeen.Before(cutoff) {
			delete(r.entries, p)
			removed++
		}
	}
	return removed
}
