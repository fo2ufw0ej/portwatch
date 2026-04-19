package envelope

import (
	"crypto/rand"
	"encoding/hex"
)

// newID returns a random 8-byte hex string suitable for use as an envelope ID.
func newID() string {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return "00000000000000000"
	}
	return hex.EncodeToString(b)
}
