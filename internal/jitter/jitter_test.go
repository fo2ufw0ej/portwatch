package jitter_test

import (
	"testing"
	"time"

	"github.com/cstevenson98/portwatch/internal/jitter"
)

func TestDefaultConfig(t *testing.T) {
	cfg := jitter.DefaultConfig()
	if cfg.Factor != 0.1 {
		t.Fatalf("expected factor 0.1, got %v", cfg.Factor)
	}
}

func TestValidate_Valid(t *testing.T) {
	if err := jitter.DefaultConfig().Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_NegativeFactor(t *testing.T) {
	cfg := jitter.Config{Factor: -0.1}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative factor")
	}
}

func TestValidate_FactorAboveOne(t *testing.T) {
	cfg := jitter.Config{Factor: 1.5}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for factor > 1")
	}
}

func TestNew_InvalidConfig(t *testing.T) {
	_, err := jitter.New(jitter.Config{Factor: 2.0})
	if err == nil {
		t.Fatal("expected error from New with invalid config")
	}
}

func TestApply_InRange(t *testing.T) {
	a, err := jitter.New(jitter.Config{Factor: 0.2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	base := 100 * time.Millisecond
	for i := 0; i < 50; i++ {
		got := a.Apply(base)
		if got < base || got > base+20*time.Millisecond {
			t.Fatalf("Apply(%v) = %v out of expected range", base, got)
		}
	}
}

func TestApplyNegative_InRange(t *testing.T) {
	a, err := jitter.New(jitter.Config{Factor: 0.2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	base := 100 * time.Millisecond
	for i := 0; i < 50; i++ {
		got := a.ApplyNegative(base)
		if got < base-20*time.Millisecond || got > base {
			t.Fatalf("ApplyNegative(%v) = %v out of expected range", base, got)
		}
	}
}

func TestApply_ZeroFactor(t *testing.T) {
	a, err := jitter.New(jitter.Config{Factor: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	base := 50 * time.Millisecond
	if got := a.Apply(base); got != base {
		t.Fatalf("expected %v, got %v", base, got)
	}
}
