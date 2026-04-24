// Package ledger tracks cumulative port-event counts across daemon runs.
// It provides a persistent tally of how many times each port has been
// seen opened or closed, useful for identifying flapping ports.
package ledger

import (
	"encoding/json"
	"os"
	"sort"
	"sync"
	"time"
)

// Entry holds the cumulative event counts for a single port.
type Entry struct {
	Port        int       `json:"port"`
	OpenCount   int       `json:"open_count"`
	CloseCount  int       `json:"close_count"`
	LastSeen    time.Time `json:"last_seen"`
}

// Ledger maintains a persistent tally of port events.
type Ledger struct {
	mu      sync.Mutex
	path    string
	entries map[int]*Entry
}

// New creates a Ledger backed by the given file path.
// Existing data is loaded if the file exists.
func New(path string) (*Ledger, error) {
	l := &Ledger{path: path, entries: make(map[int]*Entry)}
	if err := l.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return l, nil
}

// RecordOpened increments the open count for the given port.
func (l *Ledger) RecordOpened(port int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	e := l.getOrCreate(port)
	e.OpenCount++
	e.LastSeen = time.Now()
}

// RecordClosed increments the close count for the given port.
func (l *Ledger) RecordClosed(port int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	e := l.getOrCreate(port)
	e.CloseCount++
	e.LastSeen = time.Now()
}

// Get returns the entry for a port, or nil if unseen.
func (l *Ledger) Get(port int) *Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	e, ok := l.entries[port]
	if !ok {
		return nil
	}
	copy := *e
	return &copy
}

// Entries returns all entries sorted by port number.
func (l *Ledger) Entries() []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]Entry, 0, len(l.entries))
	for _, e := range l.entries {
		out = append(out, *e)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Port < out[j].Port })
	return out
}

// Save persists the ledger to disk.
func (l *Ledger) Save() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.flush()
}

func (l *Ledger) getOrCreate(port int) *Entry {
	if e, ok := l.entries[port]; ok {
		return e
	}
	e := &Entry{Port: port}
	l.entries[port] = e
	return e
}

func (l *Ledger) load() error {
	data, err := os.ReadFile(l.path)
	if err != nil {
		return err
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return err
	}
	for i := range entries {
		e := entries[i]
		l.entries[e.Port] = &e
	}
	return nil
}

func (l *Ledger) flush() error {
	entries := make([]Entry, 0, len(l.entries))
	for _, e := range l.entries {
		entries = append(entries, *e)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Port < entries[j].Port })
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(l.path, data, 0o644)
}
