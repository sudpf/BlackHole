package main

import (
	"BlackHole/internal/voidengine/app"
	"BlackHole/pkg/config"
	"BlackHole/pkg/logger"
	"flag"

	log "github.com/sirupsen/logrus"
)

func main() {
	configFile := flag.String("config-file", "voidengine.toml", "config file")
	flag.Parse()

	cfg, err := config.LoadVoidEngineConfig(*configFile)
	if err != nil {
		log.WithError(err).Fatal("Parse config file error")
	}

	logger.InitLog(cfg.LogLevel(), cfg.AppLogFile())

	log.Info(cfg.String())

	if err := app.Run(cfg); err != nil {
		log.WithError(err).Fatal("Run VoidEngine error")
	}
}
