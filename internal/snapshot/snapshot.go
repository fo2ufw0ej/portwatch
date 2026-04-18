// Package snapshot provides functionality for capturing and comparing
// point-in-time views of open ports, enabling change detection between scans.
package snapshot

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Snapshot represents a point-in-time capture of open ports on the host.
type Snapshot struct {
	Timestamp time.Time `json:"timestamp"`
	Ports     []int     `json:"ports"`
	Host      string    `json:"host"`
}

// New creates a new Snapshot with the given host and port list.
// Ports are sorted in ascending order for consistent comparison.
func New(host string, ports []int) *Snapshot {
	sorted := make([]int, len(ports))
	copy(sorted, ports)
	sort.Ints(sorted)
	return &Snapshot{
		Timestamp: time.Now(),
		Ports:     sorted,
		Host:      host,
	}
}

// Equal returns true if two snapshots contain the same set of ports.
func (s *Snapshot) Equal(other *Snapshot) bool {
	if len(s.Ports) != len(other.Ports) {
		return false
	}
	for i, p := range s.Ports {
		if other.Ports[i] != p {
			return false
		}
	}
	return true
}

// Contains returns true if the given port exists in the snapshot.
func (s *Snapshot) Contains(port int) bool {
	for _, p := range s.Ports {
		if p == port {
			return true
		}
	}
	return false
}

// Diff returns the ports opened and closed between s and a newer snapshot.
// opened contains ports present in other but not in s.
// closed contains ports present in s but not in other.
func (s *Snapshot) Diff(other *Snapshot) (opened, closed []int) {
	for _, p := range other.Ports {
		if !s.Contains(p) {
			opened = append(opened, p)
		}
	}
	for _, p := range s.Ports {
		if !other.Contains(p) {
			closed = append(closed, p)
		}
	}
	return opened, closed
}

// Summary returns a human-readable one-line description of the snapshot.
func (s *Snapshot) Summary() string {
	if len(s.Ports) == 0 {
		return fmt.Sprintf("[%s] %s — no open ports", s.Timestamp.Format(time.RFC3339), s.Host)
	}
	parts := make([]string, len(s.Ports))
	for i, p := range s.Ports {
		parts[i] = fmt.Sprintf("%d", p)
	}
	return fmt.Sprintf("[%s] %s — %d open port(s): %s",
		s.Timestamp.Format(time.RFC3339),
		s.Host,
		len(s.Ports),
		strings.Join(parts, ", "),
	)
}
