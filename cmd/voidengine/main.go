package main

import (
	"BlackHole/api/voidengine/openapi"
	"BlackHole/pkg/config"
	"BlackHole/pkg/logger"
	"flag"

	log "github.com/sirupsen/logrus"
)

func main() {
	configFile := flag.String("config-file", "./conf/BlackHole.conf", "config file")
	flag.Parse()

	err := config.ParseConfig(*configFile)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Parse config file error")
	}

	logger.InitLog(config.GlobalConfig.Log.Level, config.GlobalConfig.Log.Output)

	log.Info(config.GetConfig())

	openapi.InitApi()
	openapi.Run()
}
