// Package scorecard provides a per-port health scoring surface for portwatch.
//
// It aggregates stability scores (produced by the prestige package or any
// float64 source) into a Card that can be queried and rendered as a
// human-readable table.
//
// Typical usage:
//
//	card := scorecard.New()
//	card.Set(80, 0.95, "http")
//	card.Set(443, 0.40, "https")
//	card.Write(os.Stdout)
package scorecard
