package main

import (
	"BlackHole/internal/stash/service"
	"BlackHole/pkg/config"
	"BlackHole/pkg/logger"
	"flag"

	log "github.com/sirupsen/logrus"
)

var configFile = flag.String("f", "stash.yaml", "Specify the config file")

func main() {
	flag.Parse()

	config.ParseStashConfig(*configFile)

	logger.InitLog(config.GetStashConfig().LogLevel(), config.GetStashConfig().AppLogFile())

	log.Info(config.GetStashConfig().String())

	service.Init()
	service.Run()
}
