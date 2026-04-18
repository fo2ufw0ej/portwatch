// Package fence implements an allowlist/denylist guard for port alerts.
//
// A Guard is constructed from a Config that may specify:
//
//   - Allowlist – only ports in this list will be allowed through.
//   - Denylist  – ports in this list are always blocked, even if they
//     appear in the allowlist.
//
// When the allowlist is empty every port is allowed unless it is denied.
// This makes the zero-value Config fully permissive, consistent with the
// rest of the portwatch defaults.
package fence
