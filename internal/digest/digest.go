// Package digest computes and compares port-state digests to detect changes
// without storing full snapshots.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"time"
)

// Digest holds a hash of a port set and the time it was computed.
type Digest struct {
	Hash      string    `json:"hash"`
	ComputedAt time.Time `json:"computed_at"`
	PortCount  int       `json:"port_count"`
}

// Compute returns a Digest for the given list of ports.
func Compute(ports []int) Digest {
	sorted := make([]int, len(ports))
	copy(sorted, ports)
	sort.Ints(sorted)

	h := sha256.New()
	for _, p := range sorted {
		fmt.Fprintf(h, "%d\n", p)
	}

	return Digest{
		Hash:       hex.EncodeToString(h.Sum(nil)),
		ComputedAt: time.Now().UTC(),
		PortCount:  len(ports),
	}
}

// Equal returns true if two digests represent the same port set.
func Equal(a, b Digest) bool {
	return a.Hash == b.Hash
}

// Changed returns true if the port set has changed since the previous digest.
func Changed(prev, next Digest) bool {
	return !Equal(prev, next)
}
