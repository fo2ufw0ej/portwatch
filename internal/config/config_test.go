package config

import (
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.ScanInterval != 30*time.Second {
		t.Errorf("expected 30s, got %v", cfg.ScanInterval)
	}
	if cfg.PortRange.From != 1 || cfg.PortRange.To != 65535 {
		t.Errorf("unexpected default port range: %+v", cfg.PortRange)
	}
	if cfg.AlertOutput != "stdout" {
		t.Errorf("expected stdout, got %q", cfg.AlertOutput)
	}
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTempConfig(t, `
scan_interval: "10s"
port_range:
  from: 1024
  to: 9999
alert_output: stdout
log_level: debug
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ScanInterval != 10*time.Second {
		t.Errorf("expected 10s, got %v", cfg.ScanInterval)
	}
	if cfg.PortRange.From != 1024 || cfg.PortRange.To != 9999 {
		t.Errorf("unexpected port range: %+v", cfg.PortRange)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("expected debug, got %q", cfg.LogLevel)
	}
}

func TestLoad_InvalidInterval(t *testing.T) {
	path := writeTempConfig(t, `scan_interval: "notaduration"`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid scan_interval")
	}
}

func TestLoad_InvalidPortRange(t *testing.T) {
	path := writeTempConfig(t, `
port_range:
  from: 9000
  to: 1000
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error when from > to")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
