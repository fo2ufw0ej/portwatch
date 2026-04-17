// Package rollup coalesces rapid successive diffs into a single batched
// alert, reducing noise when many ports change in a short window.
package rollup

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Handler is called with the merged diff once the window closes.
type Handler func(diff scanner.Diff)

// Rollup buffers diffs and flushes them after a quiet window.
type Rollup struct {
	mu      sync.Mutex
	window  time.Duration
	handle  Handler
	pending *scanner.Diff
	timer   *time.Timer
}

// New creates a Rollup that waits window duration after the last Push before
// invoking handle with the merged diff.
func New(window time.Duration, handle Handler) *Rollup {
	return &Rollup{window: window, handle: handle}
}

// Push adds a diff to the pending batch, resetting the flush timer.
func (r *Rollup) Push(d scanner.Diff) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.pending == nil {
		copy := d
		r.pending = &copy
	} else {
		r.pending.Opened = union(r.pending.Opened, d.Opened)
		r.pending.Closed = union(r.pending.Closed, d.Closed)
	}

	if r.timer != nil {
		r.timer.Stop()
	}
	r.timer = time.AfterFunc(r.window, r.flush)
}

// Flush forces an immediate flush of any pending diff.
func (r *Rollup) Flush() {
	if r.timer != nil {
		r.timer.Stop()
	}
	r.flush()
}

func (r *Rollup) flush() {
	r.mu.Lock()
	d := r.pending
	r.pending = nil
	r.timer = nil
	r.mu.Unlock()

	if d != nil && (len(d.Opened) > 0 || len(d.Closed) > 0) {
		r.handle(*d)
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
