// Package checkpoint provides periodic state snapshotting with versioned
// checkpoints so the daemon can resume from a known-good scan baseline after
// a restart or crash.
package checkpoint

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry is a single versioned checkpoint written to disk.
type Entry struct {
	Version   uint64    `json:"version"`
	Ports     []int     `json:"ports"`
	CreatedAt time.Time `json:"created_at"`
}

// Store manages reading and writing checkpoint files.
type Store struct {
	path string
}

// New returns a Store that persists checkpoints at path.
func New(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("checkpoint: path must not be empty")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("checkpoint: mkdir: %w", err)
	}
	return &Store{path: path}, nil
}

// Save writes a new checkpoint, incrementing the version from the previous one.
func (s *Store) Save(ports []int) (*Entry, error) {
	var prev Entry
	if existing, err := s.Load(); err == nil {
		prev = *existing
	}

	e := &Entry{
		Version:   prev.Version + 1,
		Ports:     ports,
		CreatedAt: time.Now().UTC(),
	}

	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("checkpoint: marshal: %w", err)
	}
	if err := os.WriteFile(s.path, data, 0o644); err != nil {
		return nil, fmt.Errorf("checkpoint: write: %w", err)
	}
	return e, nil
}

// Load reads the latest checkpoint from disk.
// Returns nil, nil when no checkpoint exists yet.
func (s *Store) Load() (*Entry, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("checkpoint: read: %w", err)
	}
	var e Entry
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, fmt.Errorf("checkpoint: unmarshal: %w", err)
	}
	return &e, nil
}

// Clear removes the checkpoint file from disk.
func (s *Store) Clear() error {
	if err := os.Remove(s.path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("checkpoint: remove: %w", err)
	}
	return nil
}
