package state

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Snapshot holds the last known set of open ports and when they were recorded.
type Snapshot struct {
	Ports     []int     `json:"ports"`
	RecordedAt time.Time `json:"recorded_at"`
}

// Store persists and retrieves port snapshots to/from disk.
type Store struct {
	mu       sync.RWMutex
	filePath string
}

// NewStore creates a Store backed by the given file path.
func NewStore(filePath string) *Store {
	return &Store{filePath: filePath}
}

// Save writes the snapshot atomically to disk.
func (s *Store) Save(snap Snapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}

	tmp := s.filePath + ".tmp"
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, s.filePath)
}

// Load reads the last snapshot from disk.
// Returns an empty Snapshot (no error) when the file does not yet exist.
func (s *Store) Load() (Snapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.filePath)
	if os.IsNotExist(err) {
		return Snapshot{}, nil
	}
	if err != nil {
		return Snapshot{}, err
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, err
	}
	return snap, nil
}
