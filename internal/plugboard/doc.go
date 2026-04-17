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
package plugboard
