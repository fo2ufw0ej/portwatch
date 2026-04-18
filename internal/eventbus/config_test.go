package eventbus_test

import (
	"testing"

	"github.com/dmokel/portwatch/internal/eventbus"
)

func TestDefaultConfig(t *testing.T) {
	c := eventbus.DefaultConfig()
	if c.BufferSize != 0 {
		t.Fatalf("expected BufferSize 0, got %d", c.BufferSize)
	}
}

func TestValidate_Valid(t *testing.T) {
	c := eventbus.DefaultConfig()
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_NegativeBuffer(t *testing.T) {
	c := eventbus.Config{BufferSize: -1}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for negative BufferSize")
	}
}
