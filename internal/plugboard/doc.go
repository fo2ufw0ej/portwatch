// Package plugboard implements a lightweight publish/subscribe event bus used
// by portwatch to decouple internal components.
//
// Components that produce port-change events (e.g. the daemon scan loop) publish
// to well-known event names such as "ports.opened" and "ports.closed". Any
// number of consumers — alerters, auditors, metric recorders — subscribe to
// those events without the producer needing to know about them.
//
// Event names are arbitrary strings; by convention portwatch uses dot-separated
// namespaces (e.g. "ports.opened", "scan.error").
//
// # Well-known event names
//
// The following event names are used by portwatch core components:
//
//   - "ports.opened"  — one or more ports were found open that were not open
//     in the previous scan cycle.
//   - "ports.closed"  — one or more ports that were previously open are no
//     longer reachable.
//   - "scan.started"  — a new scan cycle has begun.
//   - "scan.finished" — a scan cycle completed successfully.
//   - "scan.error"    — a scan cycle encountered a non-fatal error; the event
//     payload contains the error details.
package plugboard
