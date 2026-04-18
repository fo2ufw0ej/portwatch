package ticker

import (
	"testing"
	"time"
)

func TestTicker_FiresAtInterval(t *testing.T) {
	tk, err := New(Config{Interval: 50 * time.Millisecond, Jitter: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer tk.Stop()

	select {
	case <-tk.C():
		// success
	case <-time.After(200 * time.Millisecond):
		t.Error("ticker did not fire within timeout")
	}
}

func TestTicker_FiresMultipleTimes(t *testing.T) {
	tk, err := New(Config{Interval: 30 * time.Millisecond, Jitter: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer tk.Stop()

	count := 0
	timeout := time.After(200 * time.Millisecond)
	for count < 3 {
		select {
		case <-tk.C():
			count++
		case <-timeout:
			t.Fatalf("only received %d ticks, expected at least 3", count)
		}
	}
}

func TestTicker_StopPreventsMoreTicks(t *testing.T) {
	tk, err := New(Config{Interval: 20 * time.Millisecond, Jitter: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// drain one tick
	select {
	case <-tk.C():
	case <-time.After(100 * time.Millisecond):
		t.Fatal("no initial tick")
	}

	tk.Stop()

	// drain channel
	for len(tk.C()) > 0 {
		<-tk.C()
	}

	select {
	case <-tk.C():
		t.Error("received tick after Stop")
	case <-time.After(80 * time.Millisecond):
		// success
	}
}

func TestTicker_WithJitter_StillFires(t *testing.T) {
	tk, err := New(Config{Interval: 40 * time.Millisecond, Jitter: 0.5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer tk.Stop()

	select {
	case <-tk.C():
		// success
	case <-time.After(300 * time.Millisecond):
		t.Error("ticker with jitter did not fire within timeout")
	}
}
