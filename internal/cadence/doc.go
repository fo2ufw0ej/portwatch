// Package cadence provides a lightweight interval-drift detector for
// portwatch's scan loop.
//
// A Tracker records successive wall-clock observations and reports
// whether the gap between them matches the configured expected period
// within a tolerance band. This makes it easy to surface situations
// where the daemon's tick rate has slowed down unexpectedly — for
// example due to system load or a blocked goroutine.
//
// Usage:
//
//	tr, _ := cadence.New(cadence.DefaultConfig())
//	tr.Observe(time.Now())
//	// … later …
//	tr.Observe(time.Now())
//	if !tr.OnTime() {
//	    log.Println("scan cadence has drifted")
//	}
package cadence
