package main

import (
	"BlackHole/api/voidengine/openapi"
	v1handler "BlackHole/api/voidengine/openapi/v1/handler"
	v1router "BlackHole/api/voidengine/openapi/v1/router"
	"BlackHole/internal/voidengine/model"
	"BlackHole/internal/voidengine/service"
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
	models, err := model.New(cfg.Database)
	if err != nil {
		log.WithError(err).Fatal("Initialize models error")
	}
	services := service.New(models)
	handlers := v1handler.New(services)

	v1router.RegisterRoutes(handlers)
	openapi.InitApi(listenAddress)
	if err := openapi.Run(listenAddress); err != nil {
		log.WithError(err).Fatal("Run OpenAPI error")
	}
}
