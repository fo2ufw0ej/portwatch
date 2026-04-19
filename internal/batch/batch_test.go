package batch_test

import (
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/batch"
	"github.com/user/portwatch/internal/scanner"
)

func defaultCfg(window time.Duration) batch.Config {
	return batch.Config{Window: window}
}

func TestPush_FlushesAfterWindow(t *testing.T) {
	var mu sync.Mutex
	var got scanner.Diff

	b, err := batch.New(defaultCfg(40*time.Millisecond), func(d scanner.Diff) {
		mu.Lock()
		got = d
		mu.Unlock()
	})
	if err != nil {
		t.Fatal(err)
	}

	b.Push(scanner.Diff{Opened: []int{80, 443}})
	time.Sleep(80 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(got.Opened) != 2 {
		t.Fatalf("expected 2 opened ports, got %d", len(got.Opened))
	}
}

func TestPush_MergesDiffs(t *testing.T) {
	var mu sync.Mutex
	var got scanner.Diff

	b, _ := batch.New(defaultCfg(60*time.Millisecond), func(d scanner.Diff) {
		mu.Lock()
		got = d
		mu.Unlock()
	})

	b.Push(scanner.Diff{Opened: []int{80}})
	b.Push(scanner.Diff{Opened: []int{443}, Closed: []int{22}})
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	sort.Ints(got.Opened)
	if len(got.Opened) != 2 || got.Opened[0] != 80 || got.Opened[1] != 443 {
		t.Fatalf("unexpected opened: %v", got.Opened)
	}
	if len(got.Closed) != 1 || got.Closed[0] != 22 {
		t.Fatalf("unexpected closed: %v", got.Closed)
	}
}

func TestFlush_Immediate(t *testing.T) {
	called := make(chan scanner.Diff, 1)
	b, _ := batch.New(defaultCfg(5*time.Second), func(d scanner.Diff) {
		called <- d
	})

	b.Push(scanner.Diff{Opened: []int{8080}})
	b.Flush()

	select {
	case d := <-called:
		if len(d.Opened) != 1 {
			t.Fatalf("expected 1 opened port, got %d", len(d.Opened))
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("flush did not fire")
	}
}

func TestFlush_EmptyDiff_NoCallback(t *testing.T) {
	called := false
	b, _ := batch.New(defaultCfg(5*time.Second), func(_ scanner.Diff) { called = true })
	b.Flush()
	if called {
		t.Fatal("handler should not be called for empty diff")
	}
}

func TestNew_InvalidConfig(t *testing.T) {
	_, err := batch.New(batch.Config{Window: 0}, func(_ scanner.Diff) {})
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}
