package watcher

import (
	"fmt"
	"net"
	"time"
)

// PortStatus represents the liveness of a specific port.
type PortStatus struct {
	Port      int
	Open      bool
	Latency   time.Duration
	CheckedAt time.Time
}

// Probe checks whether a single TCP port on the given host is open.
// It returns a PortStatus with timing information.
func Probe(host string, port int, timeout time.Duration) PortStatus {
	start := time.Now()
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", addr, timeout)
	latency := time.Since(start)

	status := PortStatus{
		Port:      port,
		Latency:   latency,
		CheckedAt: time.Now(),
	}

	if err == nil {
		conn.Close()
		status.Open = true
	}

	return status
}

// ProbeAll checks a slice of ports on the given host concurrently and
// returns a map of port number to PortStatus.
func ProbeAll(host string, ports []int, timeout time.Duration) map[int]PortStatus {
	type result struct {
		port   int
		status PortStatus
	}

	ch := make(chan result, len(ports))

	for _, p := range ports {
		go func(port int) {
			ch <- result{port: port, status: Probe(host, port, timeout)}
		}(p)
	}

	out := make(map[int]PortStatus, len(ports))
	for range ports {
		r := <-ch
		out[r.port] = r.status
	}
	return out
}
