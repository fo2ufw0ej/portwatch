// Package fingerprint provides a lightweight, order-independent hashing
// mechanism for sets of TCP port numbers.
//
// A Fingerprint is a hex-encoded SHA-256 digest computed over the sorted
// port list. It allows the daemon to quickly determine whether the set of
// open ports has changed between two consecutive scans without persisting
// or comparing full port slices.
//
// Typical usage:
//
//	prev := fingerprint.Compute(lastPorts)
//	if fingerprint.Changed(prev, currentPorts) {
//		// handle change
//	}
package fingerprint
