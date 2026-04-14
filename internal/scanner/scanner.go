package scanner

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

// PortState represents the state of a scanned port.
type PortState struct {
	Port     int
	Protocol string
	Open     bool
}

// Scanner scans local ports for open/closed state.
type Scanner struct {
	Timeout time.Duration
}

// NewScanner creates a Scanner with the given timeout.
func NewScanner(timeout time.Duration) *Scanner {
	return &Scanner{Timeout: timeout}
}

// ScanPort checks whether a single TCP port is open on localhost.
func (s *Scanner) ScanPort(port int) PortState {
	address := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", address, s.Timeout)
	if err != nil {
		return PortState{Port: port, Protocol: "tcp", Open: false}
	}
	conn.Close()
	return PortState{Port: port, Protocol: "tcp", Open: true}
}

// ScanRange scans all ports in [start, end] inclusive and returns open ones.
func (s *Scanner) ScanRange(start, end int) ([]PortState, error) {
	if start < 1 || end > 65535 || start > end {
		return nil, fmt.Errorf("invalid port range: %d-%d", start, end)
	}

	results := make([]PortState, 0)
	for port := start; port <= end; port++ {
		state := s.ScanPort(port)
		if state.Open {
			results = append(results, state)
		}
	}
	return results, nil
}

// OpenPorts returns a set of open port numbers from a slice of PortState.
func OpenPorts(states []PortState) map[int]struct{} {
	set := make(map[int]struct{}, len(states))
	for _, s := range states {
		if s.Open {
			set[s.Port] = struct{}{}
		}
	}
	return set
}
