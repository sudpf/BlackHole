package main

import (
	"BlackHole/internal/voidengine/app"
	"BlackHole/internal/voidengine/config"
	"BlackHole/pkg/logger"
	"flag"

	log "github.com/sirupsen/logrus"
)

func main() {
	configFile := flag.String("config-file", "voidengine.toml", "config file")
	flag.Parse()

	cfg, err := config.Load(*configFile)
	if err != nil {
		log.WithError(err).Fatal("Parse config file error")
	}

	if err := logger.InitLog(cfg.LogLevel(), cfg.AppLogFile(), cfg.LogSize()); err != nil {
		log.WithError(err).Fatal("Init log error")
	}

	log.Info(cfg.String())

	if err := app.Run(cfg); err != nil {
		log.WithError(err).Fatal("Run VoidEngine error")
	}
}
