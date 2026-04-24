// Package fingerprint generates a stable identity string for a set of open
// ports. It is used to detect whether the port landscape has changed between
// two consecutive scans without storing the full port list.
package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// Fingerprint is a short, stable string that uniquely identifies a port set.
type Fingerprint string

// Empty is the fingerprint of an empty port set.
const Empty Fingerprint = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

// Compute returns a deterministic SHA-256 fingerprint for the given port list.
// Port order does not matter; the list is sorted before hashing.
func Compute(ports []int) Fingerprint {
	if len(ports) == 0 {
		return Empty
	}

	sorted := make([]int, len(ports))
	copy(sorted, ports)
	sort.Ints(sorted)

	parts := make([]string, len(sorted))
	for i, p := range sorted {
		parts[i] = fmt.Sprintf("%d", p)
	}

	h := sha256.Sum256([]byte(strings.Join(parts, ",")))
	return Fingerprint(hex.EncodeToString(h[:]))
}

// Equal reports whether two fingerprints are identical.
func Equal(a, b Fingerprint) bool {
	return a == b
}

// Changed reports whether the fingerprint for newPorts differs from prev.
func Changed(prev Fingerprint, newPorts []int) bool {
	return !Equal(prev, Compute(newPorts))
}

// Short returns the first 12 hex characters of the fingerprint, suitable for
// display in log lines and audit entries.
func (f Fingerprint) Short() string {
	if len(f) <= 12 {
		return string(f)
	}
	return string(f[:12])
}

// String implements fmt.Stringer.
func (f Fingerprint) String() string {
	return string(f)
}
