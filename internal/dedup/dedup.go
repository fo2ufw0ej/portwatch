// Package dedup suppresses duplicate alerts by tracking a hash of the
// last-seen diff. An alert is only forwarded when the diff content changes.
package dedup

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Store tracks the last hash seen per named channel.
type Store struct {
	mu   sync.Mutex
	last map[string]string
}

// New returns an initialised Store.
func New() *Store {
	return &Store{last: make(map[string]string)}
}

// Changed returns true when the diff is non-empty AND its content differs from
// the previously recorded hash for key. It updates the stored hash on true.
func (s *Store) Changed(key string, d scanner.Diff) bool {
	if len(d.Opened) == 0 && len(d.Closed) == 0 {
		return false
	}
	h := hash(d)
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.last[key] == h {
		return false
	}
	s.last[key] = h
	return true
}

// Reset clears the stored hash for key.
func (s *Store) Reset(key string) {
	s.mu.Lock()
	delete(s.last, key)
	s.mu.Unlock()
}

func hash(d scanner.Diff) string {
	ports := make([]int, 0, len(d.Opened)+len(d.Closed))
	for _, p := range d.Opened {
		ports = append(ports, p*2)
	}
	for _, p := range d.Closed {
		ports = append(ports, p*2+1)
	}
	sort.Ints(ports)
	h := sha256.Sum256([]byte(fmt.Sprint(ports)))
	return fmt.Sprintf("%x", h)
}
