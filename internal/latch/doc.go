// Package latch implements a per-port sticky boolean gate.
//
// A Latch trips when an unexpected port event is observed and remains
// tripped until explicitly reset. This allows the daemon to distinguish
// between "has ever been unexpected" and "is currently unexpected",
// which is useful for generating one-shot alerts or audit entries that
// fire exactly once per monitoring epoch regardless of how many scans
// confirm the same anomaly.
//
// Usage:
//
//	l := latch.New()
//	if l.Trip(8080) {
//		// first time we've seen port 8080 in an unexpected state
//	}
//	// later, after the operator acknowledges the alert:
//	l.Reset(8080)
package latch
