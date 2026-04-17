package rollup

import (
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func collectDiff(t *testing.T, window time.Duration) (*Rollup, *[]scanner.Diff, *sync.Mutex) {
	t.Helper()
	var mu sync.Mutex
	var got []scanner.Diff
	r := New(window, func(d scanner.Diff) {
		mu.Lock()
		got = append(got, d)
		mu.Unlock()
	})
	return r, &got, &mu
}

func TestPush_FlushesAfterWindow(t *testing.T) {
	r, got, mu := collectDiff(t, 50*time.Millisecond)
	r.Push(scanner.Diff{Opened: []int{80}, Closed: []int{}})
	time.Sleep(120 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	if len(*got) != 1 {
		t.Fatalf("expected 1 flush, got %d", len(*got))
	}
}

func TestPush_MergesDiffs(t *testing.T) {
	r, got, mu := collectDiff(t, 60*time.Millisecond)
	r.Push(scanner.Diff{Opened: []int{80}, Closed: []int{22}})
	r.Push(scanner.Diff{Opened: []int{443}, Closed: []int{8080}})
	time.Sleep(150 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	if len(*got) != 1 {
		t.Fatalf("expected 1 merged flush, got %d", len(*got))
	}
	sort.Ints((*got)[0].Opened)
	if len((*got)[0].Opened) != 2 {
		t.Errorf("expected 2 opened ports, got %v", (*got)[0].Opened)
	}
}

func TestFlush_ImmediateFlush(t *testing.T) {
	r, got, mu := collectDiff(t, 5*time.Second)
	r.Push(scanner.Diff{Opened: []int{9000}, Closed: []int{}})
	r.Flush()
	time.Sleep(20 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	if len(*got) != 1 {
		t.Fatalf("expected immediate flush, got %d", len(*got))
	}
}

func TestFlush_EmptyDiff_NoCallback(t *testing.T) {
	called := false
	r := New(50*time.Millisecond, func(d scanner.Diff) { called = true })
	r.Flush()
	time.Sleep(20 * time.Millisecond)
	if called {
		t.Error("handler should not be called for empty diff")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Window <= 0 {
		t.Error("default window should be positive")
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("default config invalid: %v", err)
	}
}

func TestNewFromConfig_InvalidWindow(t *testing.T) {
	_, err := NewFromConfig(Config{Window: 0}, func(_ scanner.Diff) {})
	if err == nil {
		t.Error("expected error for zero window")
	}
}
