package main

import (
	"BlackHole/api/openapi"
	_ "BlackHole/api/openapi/v1/voidengine"
	"BlackHole/internal/model"
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
	}

	logger.InitLog(config.GetVoidEngineConfig().LogLevel(), config.GetVoidEngineConfig().AppLogFile())

	log.Info(config.GetVoidEngineConfig().String())

	openapi.InitApi()
	model.InitDB()
	openapi.Run()
}
