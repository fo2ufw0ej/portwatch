package scanner

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func startTestListener(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := ln.Addr().(*net.TCPAddr).Port
 ln.Close() }
}
_Open(t *testing.T) {
	port, stop := startTestListener(t)
	defer stop()

	s := NewScanner(500 * time.Millisecond)
	state := s.ScanPort(port)

	assert.True(t, state.Open)
	assert.Equal(t, port, state.Port)
	assert.Equal(t, "tcp", state.Protocol)
}

func TestScanPort_Closed(t *testing.T) {
	s := NewScanner(200 * time.Millisecond)
	// Port 1 is almost certainly closed in test environments.
	state := s.ScanPort(1)
	assert.False(t, state.Open)
}

func TestScanRange_InvalidRange(t *testing.T) {
	s := NewScanner(200 * time.Millisecond)
	_, err := s.ScanRange(100, 50)
	assert.Error(t, err)
}

func TestOpenPorts(t *testing.T) {
	states := []PortState{
		{Port: 80, Open: true},
		{Port: 81, Open: false},
		{Port: 443, Open: true},
	}
	set := OpenPorts(states)
	assert.Contains(t, set, 80)
	assert.Contains(t, set, 443)
	assert.NotContains(t, set, 81)
}

func TestComputeDiff(t *testing.T) {
	prev := map[int]struct{}{80: {}, 443: {}}
	curr := map[int]struct{}{443: {}, 8080: {}}

	d := ComputeDiff(prev, curr)

	assert.ElementsMatch(t, []int{8080}, d.Opened)
	assert.ElementsMatch(t, []int{80}, d.Closed)
	assert.True(t, d.HasChanges())
}

func TestComputeDiff_NoChanges(t *testing.T) {
	ports := map[int]struct{}{80: {}, 443: {}}
	d := ComputeDiff(ports, ports)
	assert.False(t, d.HasChanges())
}
