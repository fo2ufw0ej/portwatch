package scorecard

import (
	"bytes"
	"strings"
	"testing"
)

func TestSet_AndGet(t *testing.T) {
	c := New()
	c.Set(80, 0.9, "http")
	e, ok := c.Get(80)
	if !ok {
		t.Fatal("expected entry for port 80")
	}
	if e.Score != 0.9 || e.Label != "http" {
		t.Fatalf("unexpected entry: %+v", e)
	}
}

func TestGet_Missing(t *testing.T) {
	c := New()
	_, ok := c.Get(9999)
	if ok {
		t.Fatal("expected no entry")
	}
}

func TestEntries_SortedByPort(t *testing.T) {
	c := New()
	c.Set(443, 0.8, "https")
	c.Set(22, 0.5, "ssh")
	c.Set(80, 0.9, "http")
	entries := c.Entries()
	if entries[0].Port != 22 || entries[1].Port != 80 || entries[2].Port != 443 {
		t.Fatalf("unexpected order: %v", entries)
	}
}

func TestLen(t *testing.T) {
	c := New()
	if c.Len() != 0 {
		t.Fatal("expected 0")
	}
	c.Set(80, 1.0, "http")
	if c.Len() != 1 {
		t.Fatal("expected 1")
	}
}

func TestWrite_ContainsPortAndLabel(t *testing.T) {
	c := New()
	c.Set(80, 0.75, "http")
	var buf bytes.Buffer
	c.Write(&buf)
	out := buf.String()
	if !strings.Contains(out, "80") || !strings.Contains(out, "http") {
		t.Fatalf("output missing expected fields: %s", out)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config invalid: %v", err)
	}
}

func TestValidate_InvalidMinScore(t *testing.T) {
	cfg := Config{MinScore: 1.5, MaxPorts: 0}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for MinScore > 1")
	}
}

func TestUnhealthy(t *testing.T) {
	cfg := Config{MinScore: 0.5}
	if cfg.Unhealthy(0.8) {
		t.Fatal("0.8 should be healthy")
	}
	if !cfg.Unhealthy(0.3) {
		t.Fatal("0.3 should be unhealthy")
	}
}
