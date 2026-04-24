// Package ledger provides a persistent, append-friendly tally of port-event
// counts across daemon restarts.
//
// Each port that has ever been seen opened or closed is tracked with
// cumulative counters. This makes it straightforward to identify ports
// that flap frequently — oscillating between open and closed states —
// which may indicate misconfigured services or transient network issues.
//
// Usage:
//
//	l, err := ledger.New("/var/lib/portwatch/ledger.json")
//	l.RecordOpened(8080)
//	l.RecordClosed(8080)
//	_ = l.Save()
package ledger
