package labelmap

import (
	"testing"
)

func TestLookup_WellKnown(t *testing.T) {
	lm := New(nil)
	if got := lm.Lookup(22); got != "ssh" {
		t.Fatalf("expected ssh, got %s", got)
	}
	if got := lm.Lookup(443); got != "https" {
		t.Fatalf("expected https, got %s", got)
	}
}

func TestLookup_Unknown_FallsBack(t *testing.T) {
	lm := New(nil)
	if got := lm.Lookup(9999); got != "port-9999" {
		t.Fatalf("expected port-9999, got %s", got)
	}
}

func TestLookup_CustomOverridesBuiltin(t *testing.T) {
	lm := New(map[int]string{22: "my-ssh"})
	if got := lm.Lookup(22); got != "my-ssh" {
		t.Fatalf("expected my-ssh, got %s", got)
	}
}

func TestLabelAll_ReturnsMappings(t *testing.T) {
	lm := New(nil)
	result := lm.LabelAll([]int{22, 80, 9999})
	if result[22] != "ssh" {
		t.Fatalf("expected ssh for 22")
	}
	if result[80] != "http" {
		t.Fatalf("expected http for 80")
	}
	if result[9999] != "port-9999" {
		t.Fatalf("expected port-9999 for 9999")
	}
}

func TestDefaultConfig(t *testing.T) {
	c := DefaultConfig()
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_InvalidPort(t *testing.T) {
	c := Config{Custom: map[int]string{0: "bad"}}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for port 0")
	}
}

func TestNewFromConfig_Valid(t *testing.T) {
	c := Config{Custom: map[int]string{8888: "myapp"}}
	lm, err := NewFromConfig(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := lm.Lookup(8888); got != "myapp" {
		t.Fatalf("expected myapp, got %s", got)
	}
}
