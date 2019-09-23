package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLog(level string, output string) {
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
