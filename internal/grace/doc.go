// Package grace provides a graceful shutdown coordinator for portwatch.
//
// A Coordinator allows multiple goroutines or components to register
// themselves as active. When Shutdown is called, the coordinator cancels
// its context and waits up to a configurable timeout for all registered
// components to call Done. This ensures clean teardown of scanners,
// alerters, and other background workers.
//
// Usage:
//
//	cfg := grace.DefaultConfig()
//	coord, _ := grace.New(cfg)
//
//	coord.Register()
//	go func() {
//		defer coord.Done()
//		// ... do work, respect coord.Context() ...
//	}()
//
//	ok := coord.Shutdown() // waits or times out
package grace
