package daemon

import (
	"context"
	"log"
	"time"

	"github.com/yourorg/portwatch/internal/alert"
	"github.com/yourorg/portwatch/internal/config"
	"github.com/yourorg/portwatch/internal/scanner"
	"github.com/yourorg/portwatch/internal/state"
)

// Daemon orchestrates periodic port scanning, diff computation, alerting,
// and state persistence.
type Daemon struct {
	cfg     *config.Config
	scanner *scanner.Scanner
	store   *state.Store
	notifier *alert.Notifier
}

// New constructs a Daemon from the provided configuration.
func New(cfg *config.Config, storePath string) (*Daemon, error) {
	sc, err := scanner.NewScanner(cfg.PortRange.Start, cfg.PortRange.End, cfg.Timeout)
	if err != nil {
		return nil, err
	}
	notifier, err := alert.NewNotifier(cfg.AlertOutput)
	if err != nil {
		return nil, err
	}
	return &Daemon{
		cfg:      cfg,
		scanner:  sc,
		store:    state.NewStore(storePath),
		notifier: notifier,
	}, nil
}

// Run starts the daemon loop, blocking until ctx is cancelled.
func (d *Daemon) Run(ctx context.Context) error {
	log.Printf("portwatch daemon started (interval=%s, ports=%d-%d)",
		d.cfg.Interval, d.cfg.PortRange.Start, d.cfg.PortRange.End)

	if err := d.tick(); err != nil {
		log.Printf("initial scan error: %v", err)
	}

	ticker := time.NewTicker(d.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := d.tick(); err != nil {
				log.Printf("scan error: %v", err)
			}
		case <-ctx.Done():
			log.Println("portwatch daemon stopped")
			return ctx.Err()
		}
	}
}

func (d *Daemon) tick() error {
	current, err := d.scanner.OpenPorts()
	if err != nil {
		return err
	}

	prev, err := d.store.Load()
	if err != nil {
		return err
	}

	diff := scanner.ComputeDiff(prev.Ports, current)
	if err := d.notifier.Notify(diff); err != nil {
		log.Printf("alert error: %v", err)
	}

	return d.store.Save(state.Snapshot{
		Ports:      current,
		RecordedAt: time.Now().UTC(),
	})
}
