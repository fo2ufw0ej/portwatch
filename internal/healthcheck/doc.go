// Package healthcheck provides a lightweight HTTP server that exposes
// liveness and metrics endpoints for the portwatch daemon.
//
// Endpoints:
//
//	/health  — returns {"status":"ok"} when the daemon is running.
//	/metrics — returns a JSON snapshot of current runtime metrics.
//
// Usage:
//
//	srv := healthcheck.New(":9090", metricsInstance)
//	if err := srv.Start(); err != nil {
//	    log.Fatal(err)
//	}
//	defer srv.Stop()
package healthcheck
