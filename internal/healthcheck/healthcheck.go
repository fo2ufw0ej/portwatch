// Package healthcheck provides a simple HTTP health endpoint for the portwatch daemon.
package healthcheck

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

// Server exposes a lightweight HTTP health endpoint.
type Server struct {
	addr    string
	metrics *metrics.Metrics
	server  *http.Server
}

// New creates a new health check Server listening on addr.
func New(addr string, m *metrics.Metrics) *Server {
	s := &Server{addr: addr, metrics: m}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/metrics", s.handleMetrics)
	s.server = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	return s
}

// Start begins serving HTTP requests. It is non-blocking.
func (s *Server) Start() error {
	go func() { _ = s.server.ListenAndServe() }()
	return nil
}

// Stop gracefully shuts down the HTTP server.
func (s *Server) Stop() error {
	return s.server.Close()
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, `{"status":"ok"}`)
}

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	snap := s.metrics.Snapshot()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(snap)
}
