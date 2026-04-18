package window

import (
	"testing"
	"time"
)

func TestValidate_Valid(t *testing.T) {
	cfg := Config{Period: time.Second, Buckets: 2}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_NegativePeriod(t *testing.T) {
	cfg := Config{Period: -time.Second, Buckets: 2}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative period")
	}
}

func TestNewFromConfig_InvalidBuckets(t *testing.T) {
	cfg := Config{Period: time.Minute, Buckets: -1}
	_, err := NewFromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for invalid buckets")
	}
}
