package main

import (
	"BlackHole/config"
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func initLog(level string, output string) {
	//set log formatter
	logFormatter := new(log.TextFormatter)
	logFormatter.TimestampFormat = "2006-01-02 15:04:05"
	logFormatter.FullTimestamp = true
	log.SetFormatter(logFormatter)

	//set default level
	log.SetLevel(log.InfoLevel)

	//set default output
	log.SetOutput(os.Stdout)

	//set level
	logLevel, err := log.ParseLevel(level)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Parse log level error")

		logLevel = log.InfoLevel
	}
	log.SetLevel(logLevel)

	//set output
	switch {
	case output == "stderr":
		log.SetOutput(os.Stderr)
	case output == "stdout":
		log.SetOutput(os.Stdout)
	default:
		logOutput := &lumberjack.Logger{
			Filename: output,
			Compress: true,
		}

		log.SetOutput(logOutput)
	}

	return
}

func main() {
	configFile := flag.String("config-file", "./conf/BlackHole.conf", "config file")
	flag.Parse()

	err := config.ParseConfig(*configFile)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Parse config file error")
	}

	initLog(config.GlobalConfig.Log.Level, config.GlobalConfig.Log.Output)

	log.Info(config.GetConfig())

	log.WithFields(log.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")
}
