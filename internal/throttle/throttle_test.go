package throttle

import (
	"testing"
	"time"
)

func TestAllow_FirstCall(t *testing.T) {
	th := New(100 * time.Millisecond)
	if !th.Allow("scan") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_TooSoon(t *testing.T) {
	base := time.Now()
	th := New(100 * time.Millisecond)
	th.SetNow(func() time.Time { return base })
	th.Allow("scan")
	if th.Allow("scan") {
		t.Fatal("expected second immediate call to be throttled")
	}
}

func TestAllow_AfterGap(t *testing.T) {
	base := time.Now()
	th := New(100 * time.Millisecond)
	th.SetNow(func() time.Time { return base })
	th.Allow("scan")
	th.SetNow(func() time.Time { return base.Add(200 * time.Millisecond) })
	if !th.Allow("scan") {
		t.Fatal("expected call after gap to be allowed")
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	base := time.Now()
	th := New(100 * time.Millisecond)
	th.SetNow(func() time.Time { return base })
	th.Allow("a")
	if !th.Allow("b") {
		t.Fatal("key 'b' should be independent of key 'a'")
	}
}

func TestReset_ClearsEntry(t *testing.T) {
	base := time.Now()
	th := New(100 * time.Millisecond)
	th.SetNow(func() time.Time { return base })
	th.Allow("scan")
	th.Reset("scan")
	if !th.Allow("scan") {
		t.Fatal("expected allow after reset")
	}
}
