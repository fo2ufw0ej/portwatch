package history

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry records a single diff event with a timestamp.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Opened    []int     `json:"opened"`
	Closed    []int     `json:"closed"`
}

// Log holds an ordered list of diff entries.
type Log struct {
	mu      sync.RWMutex
	entries []Entry
	path    string
	maxSize int
}

// New creates a Log backed by the given file path.
// maxSize limits how many entries are kept (0 = unlimited).
func New(path string, maxSize int) *Log {
	return &Log{path: path, maxSize: maxSize}
}

// Append adds a new entry to the log and persists it.
func (l *Log) Append(opened, closed []int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	e := Entry{
		Timestamp: time.Now().UTC(),
		Opened:    opened,
		Closed:    closed,
	}
	l.entries = append(l.entries, e)

	if l.maxSize > 0 && len(l.entries) > l.maxSize {
		l.entries = l.entries[len(l.entries)-l.maxSize:]
	}

	return l.save()
}

// Entries returns a copy of all recorded entries.
func (l *Log) Entries() []Entry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]Entry, len(l.entries))
	copy(out, l.entries)
	return out
}

// Load reads persisted entries from disk.
func (l *Log) Load() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	data, err := os.ReadFile(l.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &l.entries)
}

func (l *Log) save() error {
	data, err := json.MarshalIndent(l.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(l.path, data, 0o644)
}
