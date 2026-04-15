// Package suppress implements a cooldown-based deduplication layer for
// portwatch alerts. It prevents the same port event (e.g. "port 80 opened")
// from firing repeated notifications within a configurable time window.
//
// Usage:
//
//	s := suppress.New(30 * time.Second)
//	if s.Allow("port:80:opened") {
//	    notifier.Notify(diff)
//	}
//
// Keys are arbitrary strings; by convention portwatch uses the format
// "port:<number>:<event>" where event is "opened" or "closed".
package suppress
