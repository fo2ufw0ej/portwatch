// Package decay provides an exponential decay scorer keyed by string.
//
// Scores decrease toward zero over time at a configurable half-life.
// This is useful for distinguishing ports that are consistently open from
// those that appear only transiently.
//
// Example:
//
//	s, _ := decay.NewDefault()
//	s.Add("8080", 1.0)
//	score := s.Get("8080") // slightly less than 1.0 after some time
package decay
