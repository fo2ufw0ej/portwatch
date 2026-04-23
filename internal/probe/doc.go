// Package probe implements lightweight TCP reachability checks for individual
// ports. It tracks consecutive failure counts so callers can decide when a
// port has become persistently unreachable and should trigger an alert.
//
// Basic usage:
//
//	p, _ := probe.New(probe.DefaultConfig())
//	res := p.Check(ctx, 8080)
//	if !res.Reachable && p.Failures(8080) >= cfg.MaxFailures {
//		// emit alert
//	}
package probe
