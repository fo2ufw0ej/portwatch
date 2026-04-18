// Package dedup provides duplicate-alert suppression for portwatch.
//
// A Store tracks a SHA-256 fingerprint of the last scanner.Diff delivered on
// each named channel. Callers should call Changed before forwarding an alert;
// the method returns false when the diff is empty or identical to the
// previously seen diff, preventing noisy repeated notifications.
//
// Example:
//
//	store := dedup.New()
//	if store.Changed("default", diff) {
//		notifier.Notify(diff)
//	}
package dedup
