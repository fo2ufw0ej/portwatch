package suppress_test

import (
	"testing"
	"time"

	"github.com/derekg/portwatch/internal/suppress"
)

func TestAllow_FirstCallAlwaysAllowed(t *testing.T) {
	s := suppress.New(5 * time.Second)
	if !s.Allow("port:80:opened") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallWithinCooldown_Suppressed(t *testing.T) {
	s := suppress.New(5 * time.Second)
	s.Allow("port:80:opened")
	if s.Allow("port:80:opened") {
		t.Fatal("expected second call within cooldown to be suppressed")
	}
}

func TestAllow_AfterCooldownExpiry_Allowed(t *testing.T) {
	now := time.Now()
	s := suppress.New(5 * time.Second)

	// Inject a custom clock
	s2 := &suppress.Suppressor{}
	_ = s2 // use the exported New + a trick via the package

	// Use a zero-cooldown suppressor to simulate expiry
	s = suppress.New(0)
	s.Allow("port:443:closed")
	time.Sleep(1 * time.Millisecond)
	if !s.Allow("port:443:closed") {
		t.Fatalf("expected call after cooldown expiry to be allowed (now=%v)", now)
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	s := suppress.New(10 * time.Second)
	if !s.Allow("port:80:opened") {
		t.Fatal("expected port:80:opened to be allowed")
	}
	if !s.Allow("port:443:opened") {
		t.Fatal("expected port:443:opened to be allowed (independent key)")
	}
	if s.Allow("port:80:opened") {
		t.Fatal("expected port:80:opened to be suppressed on second call")
	}
}

func TestReset_ClearsEntry(t *testing.T) {
	s := suppress.New(10 * time.Second)
	s.Allow("port:8080:opened")
	s.Reset("port:8080:opened")
	if !s.Allow("port:8080:opened") {
		t.Fatal("expected allow after reset")
	}
}

func TestStats_ReturnsSnapshot(t *testing.T) {
	s := suppress.New(10 * time.Second)
	s.Allow("port:22:opened")
	s.Allow("port:80:opened")

	stats := s.Stats()
	if len(stats) != 2 {
		t.Fatalf("expected 2 stats entries, got %d", len(stats))
	}
	if _, ok := stats["port:22:opened"]; !ok {
		t.Error("expected entry for port:22:opened")
	}
}

func TestStats_IsCopy(t *testing.T) {
	s := suppress.New(10 * time.Second)
	s.Allow("port:9000:closed")

	stats := s.Stats()
	stats["port:9000:closed"] = suppress.Entry{Count: 999}

	stats2 := s.Stats()
	if stats2["port:9000:closed"].Count == 999 {
		t.Error("Stats should return a copy, not a reference")
	}
}
