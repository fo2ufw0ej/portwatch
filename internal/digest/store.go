package digest

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Store persists the last known digest to disk.
type Store struct {
	path string
}

// NewStore returns a Store backed by the given file path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Save writes the digest to disk.
func (s *Store) Save(d Digest) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(d)
}

// Load reads the last saved digest from disk.
// Returns ErrNotFound if no digest has been saved yet.
func (s *Store) Load() (Digest, error) {
	f, err := os.Open(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Digest{}, ErrNotFound
		}
		return Digest{}, err
	}
	defer f.Close()
	var d Digest
	if err := json.NewDecoder(f).Decode(&d); err != nil {
		return Digest{}, err
	}
	return d, nil
}

// ErrNotFound is returned when no digest file exists.
var ErrNotFound = errors.New("digest: no saved digest found")
