package main

import (
	"BlackHole/internal/voidengine/app"
	"BlackHole/internal/voidengine/config"
	"BlackHole/pkg/logger"
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	if err := run(); err != nil {
		log.WithError(err).Error("VoidEngine exited")
		os.Exit(1)
	}
}

func run() error {
	configFile := flag.String("config-file", "voidengine.toml", "config file")
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
		return fmt.Errorf("run VoidEngine: %w", err)
	}

	return nil
}
