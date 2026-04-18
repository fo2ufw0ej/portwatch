// Package digest provides lightweight port-set fingerprinting.
//
// A Digest is a SHA-256 hash of a sorted list of open ports. It lets the
// daemon quickly determine whether the port state has changed between scans
// without keeping full snapshot copies in memory.
//
// Usage:
//
//	prev, _ := store.Load()
//	next  := digest.Compute(openPorts)
//	if digest.Changed(prev, next) {
//		// alert and persist
//		store.Save(next)
//	}
package digest
