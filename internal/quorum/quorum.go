// Package quorum requires a port change to be observed N consecutive
// scans before it is forwarded, reducing noise from transient blips.
package quorum

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Config holds quorum settings.
type Config struct {
	// Threshold is the number of consecutive confirmations required
	// before a change is considered stable and forwarded.
	Threshold int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{Threshold: 2}
}

// Validate returns an error if the config is invalid.
func (c Config) Validate() error {
	if c.Threshold < 1 {
		return fmt.Errorf("quorum: threshold must be >= 1, got %d", c.Threshold)
	}
	return nil
}

// entry tracks consecutive observations for a single port.
type entry struct {
	count int
	kind  string // "opened" or "closed"
}

// Quorum filters scanner.Diff values, only passing through changes that
// have been observed at least Threshold times in a row.
type Quorum struct {
	mu      sync.Mutex
	cfg     Config
	counts  map[int]*entry
}

// New creates a Quorum with the given config.
func New(cfg Config) (*Quorum, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Quorum{
		cfg:    cfg,
		counts: make(map[int]*entry),
	}, nil
}

// Confirm accepts a diff, increments internal counters, and returns a
// filtered diff containing only changes that have reached the threshold.
// Ports that disappear from the diff reset their counter.
func (q *Quorum) Confirm(diff scanner.Diff) scanner.Diff {
	q.mu.Lock()
	defer q.mu.Unlock()

	seen := make(map[int]string)
	for _, p := range diff.Opened {
		seen[p] = "opened"
	}
	for _, p := range diff.Closed {
		seen[p] = "closed"
	}

	// Reset counters for ports no longer in the diff.
	for port, e := range q.counts {
		if _, ok := seen[port]; !ok {
			_ = e
			delete(q.counts, port)
		}
	}

	var out scanner.Diff
	for port, kind := range seen {
		e, ok := q.counts[port]
		if !ok || e.kind != kind {
			q.counts[port] = &entry{count: 1, kind: kind}
			continue
		}
		e.count++
		if e.count >= q.cfg.Threshold {
			delete(q.counts, port)
			if kind == "opened" {
				out.Opened = append(out.Opened, port)
			} else {
				out.Closed = append(out.Closed, port)
			}
		}
	}
	return out
}

// Reset clears all internal counters.
func (q *Quorum) Reset() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.counts = make(map[int]*entry)
}
