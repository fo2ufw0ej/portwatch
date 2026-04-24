package cooldown

import (
	"testing"
	"time"
)

func TestValidate_Valid(t *testing.T) {
	cfg := Config{Period: 5 * time.Second}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_ZeroPeriod(t *testing.T) {
	cfg := Config{Period: 0}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero period")
	}
}

func TestValidate_NegativePeriod(t *testing.T) {
	cfg := Config{Period: -time.Second}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative period")
	}
}

func TestNewFromConfig_Valid(t *testing.T) {
	cd, err := NewFromConfig(Config{Period: 5 * time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cd == nil {
		t.Fatal("expected non-nil cooldown")
	}
}

func TestNewFromConfig_InvalidConfig(t *testing.T) {
	tests := []struct {
		name   string
		period time.Duration
	}{
		{"zero period", 0},
		{"negative period", -time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cd, err := NewFromConfig(Config{Period: tt.period})
			if err == nil {
				t.Fatalf("expected error for period %v", tt.period)
			}
			if cd != nil {
				t.Fatal("expected nil cooldown on error")
			}
		})
	}
}
