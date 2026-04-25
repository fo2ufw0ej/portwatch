// Package shadow provides a lightweight in-memory shadow store that
// tracks the most recently observed set of open ports.
//
// The shadow store is used by the daemon to skip downstream processing
// (alerting, history appends, digest updates) when consecutive scans
// return identical results, reducing noise and unnecessary I/O.
//
// Usage:
//
//	store := shadow.New()
//
//	if store.Changed(currentPorts) {
//		store.Set(currentPorts)
//		// … handle diff
//	}
package shadow
