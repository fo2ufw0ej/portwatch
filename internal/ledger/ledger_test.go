package ledger_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/wiring/portwatch/internal/ledger"
)

func tempLedgerPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "ledger.json")
}

func TestRecordOpened_IncreasesCount(t *testing.T) {
	l, err := ledger.New(tempLedgerPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	l.RecordOpened(8080)
	l.RecordOpened(8080)
	e := l.Get(8080)
	if e == nil {
		t.Fatal("expected entry, got nil")
	}
	if e.OpenCount != 2 {
		t.Errorf("OpenCount = %d, want 2", e.OpenCount)
	}
	if e.CloseCount != 0 {
		t.Errorf("CloseCount = %d, want 0", e.CloseCount)
	}
}

func TestRecordClosed_IncreasesCount(t *testing.T) {
	l, _ := ledger.New(tempLedgerPath(t))
	l.RecordClosed(443)
	e := l.Get(443)
	if e.CloseCount != 1 {
		t.Errorf("CloseCount = %d, want 1", e.CloseCount)
	}
}

func TestGet_MissingPort_ReturnsNil(t *testing.T) {
	l, _ := ledger.New(tempLedgerPath(t))
	if e := l.Get(9999); e != nil {
		t.Errorf("expected nil, got %+v", e)
	}
}

func TestSave_AndLoad_Roundtrip(t *testing.T) {
	path := tempLedgerPath(t)
	l, _ := ledger.New(path)
	l.RecordOpened(80)
	l.RecordOpened(80)
	l.RecordClosed(80)
	if err := l.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	l2, err := ledger.New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	e := l2.Get(80)
	if e == nil {
		t.Fatal("entry missing after reload")
	}
	if e.OpenCount != 2 || e.CloseCount != 1 {
		t.Errorf("got open=%d close=%d, want open=2 close=1", e.OpenCount, e.CloseCount)
	}
}

func TestEntries_SortedByPort(t *testing.T) {
	l, _ := ledger.New(tempLedgerPath(t))
	l.RecordOpened(9000)
	l.RecordOpened(80)
	l.RecordOpened(443)
	entries := l.Entries()
	ports := make([]int, len(entries))
	for i, e := range entries {
		ports[i] = e.Port
	}
	for i := 1; i < len(ports); i++ {
		if ports[i] < ports[i-1] {
			t.Errorf("entries not sorted: %v", ports)
		}
	}
}

func TestLoad_InvalidJSON_ReturnsError(t *testing.T) {
	path := tempLedgerPath(t)
	_ = os.WriteFile(path, []byte("not-json"), 0o644)
	_, err := ledger.New(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestDefaultConfig(t *testing.T) {
	c := ledger.DefaultConfig()
	if c.Path == "" {
		t.Error("DefaultConfig path should not be empty")
	}
	if c.FlappingThreshold < 1 {
		t.Errorf("FlappingThreshold = %d, want >= 1", c.FlappingThreshold)
	}
}

func TestValidate_EmptyPath(t *testing.T) {
	c := ledger.DefaultConfig()
	c.Path = ""
	if err := c.Validate(); err == nil {
		t.Error("expected error for empty path")
	}
}

func TestValidate_ZeroThreshold(t *testing.T) {
	c := ledger.DefaultConfig()
	c.FlappingThreshold = 0
	if err := c.Validate(); err == nil {
		t.Error("expected error for zero threshold")
	}
}

func TestSave_ProducesValidJSON(t *testing.T) {
	path := tempLedgerPath(t)
	l, _ := ledger.New(path)
	l.RecordOpened(22)
	_ = l.Save()
	data, _ := os.ReadFile(path)
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Errorf("saved file is not valid JSON array: %v", err)
	}
}
