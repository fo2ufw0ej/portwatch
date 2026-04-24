package roster

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time { return func() time.Time { return t } }

func TestTouch_NewPort_CreatesEntry(t *testing.T) {
	r := New()
	r.Touch(8080)
	e, ok := r.Get(8080)
	if !ok {
		t.Fatal("expected entry for port 8080")
	}
	if e.SeenCount != 1 {
		t.Fatalf("expected SeenCount=1, got %d", e.SeenCount)
	}
}

func TestTouch_ExistingPort_IncrementsCount(t *testing.T) {
	r := New()
	r.Touch(443)
	r.Touch(443)
	e, _ := r.Get(443)
	if e.SeenCount != 2 {
		t.Fatalf("expected SeenCount=2, got %d", e.SeenCount)
	}
}

func TestTouch_UpdatesLastSeen(t *testing.T) {
	base := time.Unix(1_000_000, 0)
	r := New()
	r.now = fixedNow(base)
	r.Touch(22)
	r.now = fixedNow(base.Add(time.Minute))
	r.Touch(22)
	e, _ := r.Get(22)
	if !e.LastSeen.Equal(base.Add(time.Minute)) {
		t.Fatalf("unexpected LastSeen: %v", e.LastSeen)
	}
	if !e.FirstSeen.Equal(base) {
		t.Fatalf("FirstSeen should not change: %v", e.FirstSeen)
	}
}

func TestRemove_DeletesEntry(t *testing.T) {
	r := New()
	r.Touch(3306)
	r.Remove(3306)
	if _, ok := r.Get(3306); ok {
		t.Fatal("expected entry to be removed")
	}
}

func TestLen_ReflectsCount(t *testing.T) {
	r := New()
	if r.Len() != 0 {
		t.Fatal("expected empty roster")
	}
	r.Touch(80)
	r.Touch(443)
	if r.Len() != 2 {
		t.Fatalf("expected Len=2, got %d", r.Len())
	}
}

func TestEvict_RemovesStaleEntries(t *testing.T) {
	base := time.Unix(2_000_000, 0)
	r := New()
	r.now = fixedNow(base)
	r.Touch(80)
	r.Touch(443)
	// advance time so both entries are stale
	r.now = fixedNow(base.Add(10 * time.Minute))
	r.Touch(22) // fresh
	cfg := Config{MaxAge: 5 * time.Minute}
	n := r.Evict(cfg)
	if n != 2 {
		t.Fatalf("expected 2 evictions, got %d", n)
	}
	if r.Len() != 1 {
		t.Fatalf("expected 1 remaining entry, got %d", r.Len())
	}
}

func TestEvict_ZeroMaxAge_NoOp(t *testing.T) {
	r := New()
	r.Touch(9090)
	n := r.Evict(Config{MaxAge: 0})
	if n != 0 {
		t.Fatalf("expected 0 evictions, got %d", n)
	}
	if r.Len() != 1 {
		t.Fatal("entry should not have been removed")
	}
}

func TestValidate_NegativeMaxAge(t *testing.T) {
	cfg := Config{MaxAge: -time.Second}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for negative MaxAge")
	}
}
