// Package batch collects port diff events over a time window and flushes
// them as a single aggregated payload, reducing downstream noise.
package batch

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Handler is called with the merged diff when the batch flushes.
type Handler func(diff scanner.Diff)

// Batch accumulates diffs and flushes after Window elapses or Flush is called.
type Batch struct {
	mu      sync.Mutex
	cfg     Config
	pending scanner.Diff
	timer   *time.Timer
	handler Handler
}

// New creates a Batch with the given config and flush handler.
func New(cfg Config, h Handler) (*Batch, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Batch{cfg: cfg, handler: h}, nil
}

// Push merges d into the pending diff and (re)starts the flush timer.
func (b *Batch) Push(d scanner.Diff) {
	if len(d.Opened) == 0 && len(d.Closed) == 0 {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()

	b.pending.Opened = union(b.pending.Opened, d.Opened)
	b.pending.Closed = union(b.pending.Closed, d.Closed)

	if b.timer != nil {
		b.timer.Stop()
	}
	b.timer = time.AfterFunc(b.cfg.Window, b.flush)
}

// Flush immediately fires the handler with whatever is pending.
func (b *Batch) Flush() {
	b.mu.Lock()
	if b.timer != nil {
		b.timer.Stop()
		b.timer = nil
	}
	b.mu.Unlock()
	b.flush()
}

func (b *Batch) flush() {
	b.mu.Lock()
	d := b.pending
	b.pending = scanner.Diff{}
	b.mu.Unlock()

	if len(d.Opened) > 0 || len(d.Closed) > 0 {
		b.handler(d)
	}
}

func union(a, b []int) []int {
	seen := make(map[int]struct{}, len(a)+len(b))
	for _, p := range a {
		seen[p] = struct{}{}
	}
	for _, p := range b {
		seen[p] = struct{}{}
	}
	out := make([]int, 0, len(seen))
	for p := range seen {
		out = append(out, p)
	}
	return out
}
