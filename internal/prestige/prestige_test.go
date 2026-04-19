package prestige

import (
	"testing"
)

func TestGet_DefaultScore(t *testing.T) {
	tr := New(DefaultAlpha)
	if got := tr.Get(80); got != 1.0 {
		t.Fatalf("expected 1.0 for unseen port, got %v", got)
	}
}

func TestObserve_StableIncreasesScore(t *testing.T) {
	tr := New(0.5)
	// First observe: unstable, sets score to 0.
	tr.Observe(80, false)
	if tr.Get(80) != 0.0 {
		t.Fatal("expected 0.0 after first unstable observe")
	}
	// Now stable: score should move toward 1.0.
	tr.Observe(80, true)
	got := float64(tr.Get(80))
	if got <= 0.0 || got >= 1.0 {
		t.Fatalf("expected score between 0 and 1, got %v", got)
	}
}

func TestObserve_RepeatedStable_ApproachesOne(t *testing.T) {
	tr := New(0.5)
	tr.Observe(443, false) // start low
	for i := 0; i < 20; i++ {
		tr.Observe(443, true)
	}
	if tr.Get(443) < 0.99 {
		t.Fatalf("expected score near 1.0 after many stable obs, got %v", tr.Get(443))
	}
}

func TestObserve_RepeatedUnstable_ApproachesZero(t *testing.T) {
	tr := New(0.5)
	tr.Observe(22, true) // start high
	for i := 0; i < 20; i++ {
		tr.Observe(22, false)
	}
	if tr.Get(22) > 0.01 {
		t.Fatalf("expected score near 0.0 after many unstable obs, got %v", tr.Get(22))
	}
}

func TestReset_ClearsScores(t *testing.T) {
	tr := New(DefaultAlpha)
	tr.Observe(8080, false)
	tr.Reset()
	if tr.Get(8080) != 1.0 {
		t.Fatal("expected default score after reset")
	}
}

func TestDefaultConfig(t *testing.T) {
	c := DefaultConfig()
	if err := c.Validate(); err != nil {
		t.Fatalf("default config invalid: %v", err)
	}
}

func TestValidate_InvalidAlpha(t *testing.T) {
	c := Config{Alpha: 0}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for alpha=0")
	}
}

func TestNewFromConfig_Valid(t *testing.T) {
	tr, err := NewFromConfig(DefaultConfig())
	if err != nil || tr == nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
