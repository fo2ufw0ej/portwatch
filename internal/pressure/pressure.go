// Package pressure tracks system load as a normalised score between 0.0
// (idle) and 1.0 (saturated).  Callers record observations (e.g. scan
// latency, queue depth, error rate) and query whether the daemon should
// shed work or slow its scan cadence.
package pressure

import (
	"sync"
	"time"
)

// Level is a named threshold band.
type Level int

const (
	LevelNormal  Level = iota // score < 0.5
	LevelElevated             // score >= 0.5
	LevelCritical             // score >= 0.8
)

func (l Level) String() string {
	switch l {
	case LevelElevated:
		return "elevated"
	case LevelCritical:
		return "critical"
	default:
		return "normal"
	}
}

// Gauge holds a rolling pressure score.
type Gauge struct {
	mu      sync.Mutex
	score   float64
	decay   float64 // per-second exponential decay weight (0 < decay < 1)
	updated time.Time
}

// New returns a Gauge with the given half-life for score decay.
// halfLife is the duration after which an unobserved score halves.
func New(halfLife time.Duration) *Gauge {
	if halfLife <= 0 {
		halfLife = 30 * time.Second
	}
	// decay = 0.5 ^ (1 / halfLife_seconds)
	halfSec := halfLife.Seconds()
	decayPerSec := 0.5
	_ = halfSec // used implicitly via the formula below
	// pre-compute: weight applied each second = exp(ln(0.5)/halfSec)
	// We approximate with a simple formula to avoid importing math.
	// ln(0.5) ≈ -0.693147
	decayPerSec = expApprox(-0.693147 / halfSec)

	return &Gauge{
		decay:   decayPerSec,
		updated: time.Now(),
	}
}

// Observe records a new pressure sample in [0.0, 1.0].
// Values outside that range are clamped.
func (g *Gauge) Observe(sample float64) {
	sample = clampF(sample, 0, 1)
	g.mu.Lock()
	defer g.mu.Unlock()
	g.applyDecay()
	// Blend: take the higher of current score and new sample to be
	// responsive to spikes while decaying slowly during calm periods.
	if sample > g.score {
		g.score = sample
	} else {
		// Weighted average so sustained low samples gradually reduce score.
		g.score = 0.7*g.score + 0.3*sample
	}
}

// Score returns the current pressure score in [0.0, 1.0].
func (g *Gauge) Score() float64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.applyDecay()
	return g.score
}

// Level returns the named band for the current score.
func (g *Gauge) Level() Level {
	s := g.Score()
	switch {
	case s >= 0.8:
		return LevelCritical
	case s >= 0.5:
		return LevelElevated
	default:
		return LevelNormal
	}
}

// Reset zeroes the score.
func (g *Gauge) Reset() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.score = 0
	g.updated = time.Now()
}

// applyDecay must be called with g.mu held.
func (g *Gauge) applyDecay(t ...time.Time) {
	now := time.Now()
	if len(t) > 0 {
		now = t[0]
	}
	elapsed := now.Sub(g.updated).Seconds()
	if elapsed > 0 {
		// score *= decay ^ elapsed
		g.score *= powApprox(g.decay, elapsed)
		g.updated = now
	}
}

// expApprox is a simple Taylor-series approximation of e^x for small |x|.
// Sufficient for the decay constants used here.
func expApprox(x float64) float64 {
	// Use 6-term Taylor: e^x ≈ 1 + x + x²/2 + x³/6 + x⁴/24 + x⁵/120
	x2 := x * x
	x3 := x2 * x
	x4 := x3 * x
	x5 := x4 * x
	return 1 + x + x2/2 + x3/6 + x4/24 + x5/120
}

// powApprox approximates base^exp via repeated squaring on the integer part
// plus a linear interpolation for the fractional part.
func powApprox(base, exp float64) float64 {
	if exp <= 0 {
		return 1
	}
	result := 1.0
	for exp >= 1 {
		result *= base
		exp--
	}
	// fractional remainder: base^frac ≈ 1 + frac*(base-1)
	result *= 1 + exp*(base-1)
	return result
}

func clampF(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
