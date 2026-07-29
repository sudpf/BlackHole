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

	config.ParseStashConfig(*configFile)
	cfg := config.GetStashConfig()

	logger.InitLog(cfg.LogLevel(), cfg.AppLogFile())

	log.Info(cfg.String())

	if err := app.Run(cfg); err != nil {
		log.WithError(err).Fatal("Run Stash error")
	}
}
