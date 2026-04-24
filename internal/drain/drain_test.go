package drain_test

import (
	"context"
	"testing"
	"time"

	"github.com/stevezaluk/portwatch/internal/drain"
)

func defaultCfg() drain.Config {
	return drain.Config{
		Capacity:     8,
		FlushTimeout: 2 * time.Second,
	}
}

func TestNew_InvalidCapacity(t *testing.T) {
	cfg := defaultCfg()
	cfg.Capacity = 0
	_, err := drain.New[int](cfg)
	if err == nil {
		t.Fatal("expected error for zero capacity")
	}
}

func TestNew_InvalidFlushTimeout(t *testing.T) {
	cfg := defaultCfg()
	cfg.FlushTimeout = 0
	_, err := drain.New[string](cfg)
	if err == nil {
		t.Fatal("expected error for zero flush timeout")
	}
}

func TestPush_AndConsume(t *testing.T) {
	d, err := drain.New[int](defaultCfg())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()
	for i := 0; i < 5; i++ {
		if !d.Push(ctx, i) {
			t.Fatalf("push %d failed", i)
		}
	}

	go func() {
		for range d.C() {
			d.Done()
		}
	}()

	if !d.Drain() {
		t.Fatal("drain timed out")
	}
}

func TestDrain_Timeout(t *testing.T) {
	cfg := drain.Config{
		Capacity:     4,
		FlushTimeout: 50 * time.Millisecond,
	}
	d, err := drain.New[int](cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()
	d.Push(ctx, 1) // pushed but never consumed or marked done

	// Consumer reads but never calls Done.
	go func() {
		<-d.C()
	}()

	if d.Drain() {
		t.Fatal("expected drain to time out")
	}
}

func TestPush_ContextCancelled(t *testing.T) {
	cfg := drain.Config{
		Capacity:     1,
		FlushTimeout: time.Second,
	}
	d, err := drain.New[int](cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	d.Push(ctx, 1) // fills the single slot
	cancel()

	if d.Push(ctx, 2) {
		t.Fatal("expected push to fail with cancelled context")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := drain.DefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config invalid: %v", err)
	}
}
