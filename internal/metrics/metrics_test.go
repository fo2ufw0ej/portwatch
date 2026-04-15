package metrics_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

func TestNew_ZeroValues(t *testing.T) {
	c := metrics.New()
	s := c.Snapshot()
	if s.ScansTotal != 0 {
		t.Errorf("expected 0 scans, got %d", s.ScansTotal)
	}
	if s.AlertsTotal != 0 {
		t.Errorf("expected 0 alerts, got %d", s.AlertsTotal)
	}
}

func TestRecordScan(t *testing.T) {
	c := metrics.New()
	c.RecordScan(42)
	c.RecordScan(10)
	s := c.Snapshot()
	if s.ScansTotal != 2 {
		t.Errorf("expected 2 scans, got %d", s.ScansTotal)
	}
	if s.LastScanPorts != 10 {
		t.Errorf("expected last scan ports 10, got %d", s.LastScanPorts)
	}
	if s.LastScanAt.IsZero() {
		t.Error("expected LastScanAt to be set")
	}
}

func TestRecordAlert(t *testing.T) {
	c := metrics.New()
	c.RecordAlert()
	c.RecordAlert()
	c.RecordAlert()
	s := c.Snapshot()
	if s.AlertsTotal != 3 {
		t.Errorf("expected 3 alerts, got %d", s.AlertsTotal)
	}
}

func TestSnapshot_Uptime(t *testing.T) {
	c := metrics.New()
	time.Sleep(10 * time.Millisecond)
	s := c.Snapshot()
	if s.Uptime < 0 {
		t.Error("uptime should not be negative")
	}
}

func TestSnapshot_IsCopy(t *testing.T) {
	c := metrics.New()
	c.RecordScan(5)
	s1 := c.Snapshot()
	c.RecordScan(5)
	s2 := c.Snapshot()
	if s1.ScansTotal == s2.ScansTotal {
		t.Error("snapshots should be independent copies")
	}
}
