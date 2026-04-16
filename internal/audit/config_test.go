package audit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arnavsurve/portwatch/internal/audit"
)

func TestDefaultConfig(t *testing.T) {
	c := audit.DefaultConfig()
	if c.Enabled {
		t.Error("expected disabled by default")
	}
	if c.Format != "text" {
		t.Errorf("expected text format, got %s", c.Format)
	}
}

func TestValidate_ValidFormats(t *testing.T) {
	for _, f := range []string{"text", "json"} {
		c := audit.Config{Format: f}
		if err := c.Validate(); err != nil {
			t.Errorf("expected valid format %s, got error: %v", f, err)
		}
	}
}

func TestValidate_InvalidFormat(t *testing.T) {
	c := audit.Config{Format: "xml"}
	if err := c.Validate(); err == nil {
		t.Error("expected error for invalid format")
	}
}

func TestOpen_EmptyPath_ReturnsStdout(t *testing.T) {
	c := audit.Config{Format: "text", Path: ""}
	w, err := c.Open()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer w.Close()
	if w == nil {
		t.Error("expected non-nil writer")
	}
}

func TestOpen_FilePath_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "audit.log")
	c := audit.Config{Format: "json", Path: p}
	w, err := c.Open()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w.Close()
	if _, err := os.Stat(p); os.IsNotExist(err) {
		t.Error("expected audit log file to be created")
	}
}
