package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourorg/portwatch/internal/config"
	"github.com/yourorg/portwatch/internal/daemon"
)

func main() {
	cfgPath := flag.String("config", "", "path to portwatch config file (TOML)")
	stateFile := flag.String("state", "/var/lib/portwatch/state.json", "path to state file")
	flag.Parse()

	var cfg *config.Config
	var err error

	if *cfgPath != "" {
		cfg, err = config.Load(*cfgPath)
		if err != nil {
			log.Fatalf("failed to load config %q: %v", *cfgPath, err)
		}
	} else {
		cfg = config.DefaultConfig()
		log.Println("no config file specified, using defaults")
	}

	d, err := daemon.New(cfg, *stateFile)
	if err != nil {
		log.Fatalf("failed to create daemon: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := d.Run(ctx); err != nil && err != context.Canceled {
		log.Fatalf("daemon exited with error: %v", err)
	}
}
