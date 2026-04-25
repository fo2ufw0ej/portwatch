// Package clamp enforces numeric bounds on port counts or any integer
// value observed during a portwatch scan cycle.
//
// Typical usage:
//
//	cfg := clamp.DefaultConfig()
//	cfg.Max = 1024 // restrict to well-known ports
//	c, err := clamp.New(cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	v, changed := c.Clamp(2048) // returns 1024, true
package clamp
