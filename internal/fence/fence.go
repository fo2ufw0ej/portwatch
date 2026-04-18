// Package fence provides a port allowlist/denylist guard that decides
// whether a scanned port should trigger an alert.
package fence

import "fmt"

// Config holds the fence configuration.
type Config struct {
	// Allowlist, when non-empty, only permits alerts for listed ports.
	Allowlist []int
	// Denylist suppresses alerts for listed ports regardless of allowlist.
	Denylist []int
}

// DefaultConfig returns a permissive Config with no restrictions.
func DefaultConfig() Config {
	return Config{}
}

// Validate checks the config for obvious errors.
func (c Config) Validate() error {
	for _, p := range append(c.Allowlist, c.Denylist...) {
		if p < 1 || p > 65535 {
			return fmt.Errorf("fence: port %d out of range [1, 65535]", p)
		}
	}
	return nil
}

// Guard enforces the fence rules.
type Guard struct {
	allowSet map[int]struct{}
	denySet  map[int]struct{}
}

// New creates a Guard from cfg. Returns an error if cfg is invalid.
func New(cfg Config) (*Guard, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	g := &Guard{
		allowSet: toSet(cfg.Allowlist),
		denySet:  toSet(cfg.Denylist),
	}
	return g, nil
}

// Allow returns true when port p should produce an alert.
func (g *Guard) Allow(p int) bool {
	if _, denied := g.denySet[p]; denied {
		return false
	}
	if len(g.allowSet) == 0 {
		return true
	}
	_, ok := g.allowSet[p]
	return ok
}

// Filter returns only the ports from ports that pass the guard.
func (g *Guard) Filter(ports []int) []int {
	out := make([]int, 0, len(ports))
	for _, p := range ports {
		if g.Allow(p) {
			out = append(out, p)
		}
	}
	return out
}

func toSet(ports []int) map[int]struct{} {
	s := make(map[int]struct{}, len(ports))
	for _, p := range ports {
		s[p] = struct{}{}
	}
	return s
}
