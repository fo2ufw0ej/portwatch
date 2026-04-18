// Package eventbus implements a lightweight synchronous publish/subscribe
// event bus used internally by portwatch to decouple components.
//
// Producers (e.g. the scanner daemon) publish Events; consumers (e.g. the
// alerter, auditor, metrics recorder) subscribe to specific EventTypes
// without needing direct references to one another.
package eventbus
