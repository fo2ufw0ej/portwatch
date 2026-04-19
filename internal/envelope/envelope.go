// Package envelope wraps a scanner diff with metadata for downstream consumers.
package envelope

import (
	"time"

	"github.com/example/portwatch/internal/scanner"
)

// Envelope wraps a diff with contextual metadata.
type Envelope struct {
	ID        string
	Timestamp time.Time
	Hostname  string
	Diff      scanner.Diff
	Tags      map[string]string
}

// Config holds options for the envelope builder.
type Config struct {
	Hostname string
	Tags     map[string]string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Tags: make(map[string]string),
	}
}

// Builder creates Envelopes from diffs.
type Builder struct {
	cfg Config
}

// New returns a Builder using cfg.
func New(cfg Config) (*Builder, error) {
	if cfg.Tags == nil {
		cfg.Tags = make(map[string]string)
	}
	return &Builder{cfg: cfg}, nil
}

// Wrap creates an Envelope for the given diff.
func (b *Builder) Wrap(diff scanner.Diff) Envelope {
	tags := make(map[string]string, len(b.cfg.Tags))
	for k, v := range b.cfg.Tags {
		tags[k] = v
	}
	return Envelope{
		ID:        newID(),
		Timestamp: time.Now().UTC(),
		Hostname:  b.cfg.Hostname,
		Diff:      diff,
		Tags:      tags,
	}
}

// IsEmpty reports whether the envelope carries no changes.
func (e Envelope) IsEmpty() bool {
	return len(e.Diff.Opened) == 0 && len(e.Diff.Closed) == 0
}
