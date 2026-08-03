package main

import (
	"BlackHole/internal/stash/app"
	"BlackHole/internal/stash/config"
	"BlackHole/pkg/logger"
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

var configFile = flag.String("f", "stash.yaml", "Specify the config file")

func main() {
	if err := run(); err != nil {
		log.WithError(err).Error("Stash exited")
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()

	cfg, err := config.Load(*configFile)
	if err != nil {
		return fmt.Errorf("parse config file: %w", err)
	}

	if err := logger.InitLog(cfg.LogLevel(), cfg.AppLogFile(), cfg.LogSize()); err != nil {
		return fmt.Errorf("initialize log: %w", err)
	}

	log.Info(cfg.String())

	if err := app.Run(cfg); err != nil {
		return fmt.Errorf("run Stash: %w", err)
	}

	return nil
}
