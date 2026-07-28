package main

import (
	"BlackHole/api/voidengine/openapi"
	v1router "BlackHole/api/voidengine/openapi/v1/router"
	"BlackHole/internal/voidengine/model"
	"BlackHole/pkg/config"
	"BlackHole/pkg/logger"
	"flag"

	log "github.com/sirupsen/logrus"
)

func main() {
	configFile := flag.String("config-file", "voidengine.conf", "config file")
	flag.Parse()

	err := config.ParseVoidEngineConfig(*configFile)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Parse config file error")
		return
	}

	cfg := config.GetVoidEngineConfig()
	logger.InitLog(cfg.LogLevel(), cfg.AppLogFile())

	log.Info(cfg.String())

	listenAddress := cfg.ListenHTTP()
	v1router.RegisterRoutes()
	openapi.InitApi(listenAddress)
	model.InitDB(cfg.Database)
	if err := openapi.Run(listenAddress); err != nil {
		log.WithError(err).Fatal("Run OpenAPI error")
	}
}
