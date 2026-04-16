// Package throttle provides a per-key time-based throttle that prevents
// operations from executing more frequently than a configured minimum gap.
//
// Unlike a rate limiter (which tracks burst budgets), the throttle simply
// enforces a hard minimum interval between consecutive calls for the same key.
// This is useful for scan loops where back-to-back triggers should be coalesced
// without the complexity of token-bucket accounting.
//
// Example:
//
//	th := throttle.New(5 * time.Second)
//	if th.Allow("port-scan") {
//		// run scan
//	}
package throttle
