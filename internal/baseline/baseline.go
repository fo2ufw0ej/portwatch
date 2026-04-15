// Package baseline manages the trusted port baseline for portwatch.
// A baseline represents the set of ports that are considered "expected" to be
// open. Deviations from the baseline trigger alerts.
package baseline

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"time"
)

// Entry holds a saved baseline snapshot.
type Entry struct {
	Ports     []int     `json:"ports"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Store persists and retrieves the port baseline.
type Store struct {
	path string
}

// New returns a Store backed by the given file path.
func New(path string) *Store {
	return &Store{path: path}
}

// Save writes ports as the current baseline, creating or overwriting the file.
func (s *Store) Save(ports []int) error {
	if ports == nil {
		ports = []int{}
	}
	sorted := make([]int, len(ports))
	copy(sorted, ports)
	sort.Ints(sorted)

	now := time.Now().UTC()
	existing, err := s.Load()
	createdAt := now
	if err == nil {
		createdAt = existing.CreatedAt
	}

	entry := Entry{
		Ports:     sorted,
		CreatedAt: createdAt,
		UpdatedAt: now,
	}

	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

// Load reads the baseline from disk. Returns an error if the file is missing
// or malformed.
func (s *Store) Load() (*Entry, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errors.New("baseline: no baseline file found")
		}
		return nil, err
	}
	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

// Clear removes the baseline file from disk.
func (s *Store) Clear() error {
	err := os.Remove(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
