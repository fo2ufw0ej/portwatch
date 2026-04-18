// Package pipeline provides a composable, step-based processing chain for
// portwatch scan results.
//
// A Pipeline is built from one or more Steps — plain functions that accept a
// context and a slice of open port numbers and return a (possibly transformed)
// slice plus an error.
//
// Built-in helpers (ScanStep, FilterStep) adapt existing portwatch components
// into Steps so callers can mix and match behaviour without coupling packages
// directly to one another.
//
// Example:
//
//	pl := pipeline.New(
//		pipeline.ScanStep(sc),
//		pipeline.FilterStep(rule),
//	)
//	ports, err := pl.Run(ctx, nil)
package pipeline
