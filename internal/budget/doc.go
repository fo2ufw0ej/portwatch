// Package budget provides a sliding-window error budget tracker for
// portwatch.
//
// An error budget quantifies how much unreliability is acceptable over a
// given time window. Typical usage:
//
//	cfg := budget.DefaultConfig()
//	b, err := budget.New(cfg)
//	if err != nil { /* handle */ }
//
//	// Record outcomes as they happen:
//	b.Record(true)  // success
//	b.Record(false) // failure
//
//	// Check whether the budget is still available:
//	if b.Exhausted() {
//		// suppress non-critical alerts, back off, etc.
//	}
//
// The tracker uses a ring of fixed-size time buckets so that old
// observations naturally expire without requiring a background goroutine.
package budget
