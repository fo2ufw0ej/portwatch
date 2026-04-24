// Package checkpoint implements versioned on-disk checkpoints for portwatch.
//
// A checkpoint captures the set of open ports observed at a point in time
// together with a monotonically increasing version number.  On startup the
// daemon loads the most-recent checkpoint so it can compare the current scan
// against a known-good baseline rather than treating every port as newly
// opened.
//
// Usage:
//
//	store, err := checkpoint.New("/var/lib/portwatch/checkpoint.json")
//	entry, err := store.Save([]int{22, 80, 443})
//	latest, err := store.Load()
package checkpoint
