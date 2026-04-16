package debounce_test

import (
	"sync"
	"testing"
	"time"

	"github.com/clambin/portwatch/internal/debounce"
)

func TestTrigger_FiresAfterDelay(t *testing.T) {
	var mu sync.Mutex
	fired := []string{}

	d := debounce.New(50*time.Millisecond, func(key string) {
		mu.Lock()
		fired = append(fired, key)
		mu.Unlock()
	})

	d.Trigger("port-80")
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(fired) != 1 || fired[0] != "port-80" {
		t.Fatalf("expected [port-80], got %v", fired)
	}
}

func TestTrigger_Resets(t *testing.T) {
	var mu sync.Mutex
	count := 0

	d := debounce.New(60*time.Millisecond, func(_ string) {
		mu.Lock()
		count++
		mu.Unlock()
	})

	// Trigger repeatedly; should only fire once.
	for i := 0; i < 5; i++ {
		d.Trigger("port-443")
		time.Sleep(20 * time.Millisecond)
	}
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if count != 1 {
		t.Fatalf("expected 1 fire, got %d", count)
	}
}

func TestCancel_PreventsFire(t *testing.T) {
	fired := false
	d := debounce.New(80*time.Millisecond, func(_ string) {
		fired = true
	})

	d.Trigger("port-22")
	d.Cancel("port-22")
	time.Sleep(120 * time.Millisecond)

	if fired {
		t.Fatal("callback should not have fired after cancel")
	}
}

func TestPending(t *testing.T) {
	d := debounce.New(200*time.Millisecond, func(_ string) {})

	if d.Pending() != 0 {
		t.Fatal("expected 0 pending")
	}
	d.Trigger("a")
	d.Trigger("b")
	if d.Pending() != 2 {
		t.Fatalf("expected 2 pending, got %d", d.Pending())
	}
	d.Cancel("a")
	if d.Pending() != 1 {
		t.Fatalf("expected 1 pending, got %d", d.Pending())
	}
}
