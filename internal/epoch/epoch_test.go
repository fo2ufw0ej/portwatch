package epoch_test

import (
	"sync"
	"testing"
	"time"

	"github.com/danvolchek/portwatch/internal/epoch"
)

func TestNew_ZeroState(t *testing.T) {
	c := epoch.New()
	e := c.Current()
	if e.Number != 0 {
		t.Fatalf("expected initial number 0, got %d", e.Number)
	}
	if !e.StartedAt.IsZero() {
		t.Fatalf("expected zero StartedAt before first advance")
	}
}

func TestAdvance_IncrementsNumber(t *testing.T) {
	c := epoch.New()
	e1 := c.Advance()
	e2 := c.Advance()

	if e1.Number != 1 {
		t.Fatalf("expected first epoch 1, got %d", e1.Number)
	}
	if e2.Number != 2 {
		t.Fatalf("expected second epoch 2, got %d", e2.Number)
	}
}

func TestAdvance_SetsStartedAt(t *testing.T) {
	before := time.Now()
	c := epoch.New()
	e := c.Advance()
	after := time.Now()

	if e.StartedAt.Before(before) || e.StartedAt.After(after) {
		t.Fatalf("StartedAt %v not between %v and %v", e.StartedAt, before, after)
	}
}

func TestCurrent_ReflectsLastAdvance(t *testing.T) {
	c := epoch.New()
	advanced := c.Advance()
	current := c.Current()

	if current.Number != advanced.Number {
		t.Fatalf("Current().Number = %d, want %d", current.Number, advanced.Number)
	}
}

func TestReset_ClearsState(t *testing.T) {
	c := epoch.New()
	c.Advance()
	c.Advance()
	c.Reset()

	e := c.Current()
	if e.Number != 0 {
		t.Fatalf("expected 0 after reset, got %d", e.Number)
	}
	if !e.StartedAt.IsZero() {
		t.Fatalf("expected zero StartedAt after reset")
	}
}

func TestString_Format(t *testing.T) {
	c := epoch.New()
	e := c.Advance()
	s := e.String()
	if len(s) == 0 {
		t.Fatal("String() returned empty string")
	}
}

func TestAdvance_ConcurrentSafe(t *testing.T) {
	c := epoch.New()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Advance()
		}()
	}
	wg.Wait()

	e := c.Current()
	if e.Number != 50 {
		t.Fatalf("expected epoch 50 after 50 concurrent advances, got %d", e.Number)
	}
}
