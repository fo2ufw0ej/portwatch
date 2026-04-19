package envelope_test

import (
	"testing"

	"github.com/example/portwatch/internal/envelope"
	"github.com/example/portwatch/internal/scanner"
)

func TestWrap_PopulatesFields(t *testing.T) {
	cfg := envelope.DefaultConfig()
	cfg.Hostname = "testhost"
	cfg.Tags = map[string]string{"env": "test"}

	b, err := envelope.New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	diff := scanner.Diff{Opened: []int{8080}, Closed: []int{22}}
	e := b.Wrap(diff)

	if e.Hostname != "testhost" {
		t.Errorf("hostname = %q, want testhost", e.Hostname)
	}
	if e.Tags["env"] != "test" {
		t.Errorf("tag env = %q, want test", e.Tags["env"])
	}
	if e.ID == "" {
		t.Error("expected non-empty ID")
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
	if len(e.Diff.Opened) != 1 || e.Diff.Opened[0] != 8080 {
		t.Errorf("unexpected diff.Opened: %v", e.Diff.Opened)
	}
}

func TestIsEmpty_True(t *testing.T) {
	b, _ := envelope.New(envelope.DefaultConfig())
	e := b.Wrap(scanner.Diff{})
	if !e.IsEmpty() {
		t.Error("expected IsEmpty true for empty diff")
	}
}

func TestIsEmpty_False(t *testing.T) {
	b, _ := envelope.New(envelope.DefaultConfig())
	e := b.Wrap(scanner.Diff{Opened: []int{443}})
	if e.IsEmpty() {
		t.Error("expected IsEmpty false for non-empty diff")
	}
}

func TestWrap_TagsAreCopied(t *testing.T) {
	cfg := envelope.DefaultConfig()
	cfg.Tags = map[string]string{"k": "v"}
	b, _ := envelope.New(cfg)

	e := b.Wrap(scanner.Diff{})
	e.Tags["k"] = "mutated"

	e2 := b.Wrap(scanner.Diff{})
	if e2.Tags["k"] != "v" {
		t.Errorf("tags not copied; got %q", e2.Tags["k"])
	}
}

func TestWrap_UniqueIDs(t *testing.T) {
	b, _ := envelope.New(envelope.DefaultConfig())
	e1 := b.Wrap(scanner.Diff{})
	e2 := b.Wrap(scanner.Diff{})
	if e1.ID == e2.ID {
		t.Errorf("expected unique IDs, got %q twice", e1.ID)
	}
}
