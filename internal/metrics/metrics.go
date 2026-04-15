package metrics

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time summary of daemon activity.
type Snapshot struct {
	ScansTotal    int64         `json:"scans_total"`
	AlertsTotal   int64         `json:"alerts_total"`
	LastScanAt    time.Time     `json:"last_scan_at"`
	LastScanPorts int           `json:"last_scan_ports"`
	Uptime        time.Duration `json:"uptime"`
}

// Collector tracks runtime counters for portwatch.
type Collector struct {
	mu          sync.RWMutex
	scansTotal  int64
	alertsTotal int64
	lastScanAt  time.Time
	lastScanPorts int
	startedAt   time.Time
}

// New returns an initialised Collector.
func New() *Collector {
	return &Collector{startedAt: time.Now()}
}

// RecordScan updates scan-related counters.
func (c *Collector) RecordScan(portCount int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.scansTotal++
	c.lastScanAt = time.Now()
	c.lastScanPorts = portCount
}

// RecordAlert increments the alert counter.
func (c *Collector) RecordAlert() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.alertsTotal++
}

// Snapshot returns a copy of current metrics.
func (c *Collector) Snapshot() Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return Snapshot{
		ScansTotal:    c.scansTotal,
		AlertsTotal:   c.alertsTotal,
		LastScanAt:    c.lastScanAt,
		LastScanPorts: c.lastScanPorts,
		Uptime:        time.Since(c.startedAt).Truncate(time.Second),
	}
}
