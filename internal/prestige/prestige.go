// Package prestige tracks port "reputation" — how stable a port's open/closed
// state has been over recent scans. A port that flaps frequently gets a lower
// score; one that stays consistent trends toward 1.0.
package prestige

import "sync"

// Score is a value in [0.0, 1.0] representing port stability.
type Score float64

// Tracker maintains exponential-moving-average stability scores per port.
type Tracker struct	{
	mu     sync.Mutex
	scores map[int]float64
	alpha  float64 // EMA smoothing factor
}

// DefaultAlpha is a reasonable smoothing factor (higher = faster decay).
const DefaultAlpha = 0.3

// New returns a Tracker with the given EMA alpha (0 < alpha <= 1).
func New(alpha float64) *Tracker {
	if alpha <= 0 || alpha > 1 {
		alpha = DefaultAlpha
	}
	return &Tracker{
		scores: make(map[int]float64),
		alpha:  alpha,
	}
}

// Observe records whether a port was stable (unchanged) in this scan cycle.
// stable=true nudges the score up; stable=false nudges it down.
func (t *Tracker) Observe(port int, stable bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	current, ok := t.scores[port]
	if !ok {
		if stable {
			current = 1.0
		} else {
			current = 0.0
		}
		t.scores[port] = current
		return
	}
	var target float64
	if stable {
		target = 1.0
	}
	t.scores[port] = current + t.alpha*(target-current)
}

// Get returns the current stability Score for a port (default 1.0 if unseen).
func (t *Tracker) Get(port int) Score {
	t.mu.Lock()
	defer t.mu.Unlock()
	if s, ok := t.scores[port]; ok {
		return Score(s)
	}
	return Score(1.0)
}

// Reset removes all tracked scores.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.scores = make(map[int]float64)
}
