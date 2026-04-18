// Package cooldown implements a per-key quiet-period guard.
//
// It is used to prevent repeated alerts or actions for the same event
// within a configurable duration. Unlike ratelimit (which uses a sliding
// window with a burst allowance), cooldown enforces a strict minimum gap
// between successive Allow calls for the same key.
//
// Example:
//
//	cd := cooldown.New(10 * time.Second)
//	if cd.Allow("port:8080") {
//		// send alert
//	}
package cooldown
