// Package burst provides a sliding-window event counter used to detect
// short-lived traffic spikes on a per-key basis.
//
// Typical usage:
//
//	tr, err := burst.New(burst.DefaultConfig())
//	if err != nil { ... }
//
//	tr.Record("port:8080")
//	if tr.Spiking("port:8080") {
//	    // alert: too many open/close events in the window
//	}
package burst
