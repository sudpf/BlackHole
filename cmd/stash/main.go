package main

import (
	"BlackHole/internal/stash/app"
	"BlackHole/pkg/config"
	"BlackHole/pkg/logger"
	"flag"

	log "github.com/sirupsen/logrus"
)

var configFile = flag.String("f", "stash.yaml", "Specify the config file")

func main() {
	flag.Parse()

	cfg, err := config.LoadStashConfig(*configFile)
	if err != nil {
		log.WithError(err).Fatal("Parse config file error")
	}

	if err := logger.InitLog(cfg.LogLevel(), cfg.AppLogFile(), cfg.LogSize()); err != nil {
		log.WithError(err).Fatal("Init log error")
	}

	log.Info(cfg.String())

	if err := app.Run(cfg); err != nil {
		log.WithError(err).Fatal("Run Stash error")
	}
}
