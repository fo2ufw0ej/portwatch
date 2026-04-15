package ratelimit_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/ratelimit"
)

func TestAllow_UnderBurst(t *testing.T) {
	l := ratelimit.New(time.Minute, 3)
	for i := 0; i < 3; i++ {
		if !l.Allow("port:80") {
			t.Fatalf("expected Allow to return true on call %d", i+1)
		}
	}
}

func TestAllow_ExceedsBurst(t *testing.T) {
	l := ratelimit.New(time.Minute, 2)
	l.Allow("port:443")
	l.Allow("port:443")
	if l.Allow("port:443") {
		t.Fatal("expected Allow to return false after burst exceeded")
	}
}

func TestAllow_WindowExpiry(t *testing.T) {
	now := time.Now()
	l := ratelimit.New(50*time.Millisecond, 1)

	// Inject a custom clock so we can control time.
	// We reach into the struct via the exported constructor.
	_ = now

	// Use real sleep for simplicity.
	l.Allow("port:22")
	if l.Allow("port:22") {
		t.Fatal("second call within window should be denied")
	}
	time.Sleep(60 * time.Millisecond)
	if !l.Allow("port:22") {
		t.Fatal("call after window expiry should be allowed")
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	l := ratelimit.New(time.Minute, 1)
	if !l.Allow("port:80") {
		t.Fatal("first key should be allowed")
	}
	if !l.Allow("port:443") {
		t.Fatal("different key should be allowed independently")
	}
}

func TestReset(t *testing.T) {
	l := ratelimit.New(time.Minute, 1)
	l.Allow("port:8080")
	l.Reset("port:8080")
	if !l.Allow("port:8080") {
		t.Fatal("after Reset, key should be allowed again")
	}
}

func TestLen(t *testing.T) {
	l := ratelimit.New(time.Minute, 5)
	l.Allow("port:9090")
	l.Allow("port:9090")
	l.Allow("port:9090")
	if got := l.Len("port:9090"); got != 3 {
		t.Fatalf("expected Len 3, got %d", got)
	}
}
