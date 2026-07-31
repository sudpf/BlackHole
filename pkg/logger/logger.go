package logger

import (
	"BlackHole/pkg/requestctx"
	"BlackHole/pkg/units"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLog(level string, output string, size string) error {
	log.SetFormatter(JSONFormatter())

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
		output, err := RotatingWriter(output, size)
		if err != nil {
			return err
		}
		log.SetOutput(output)
	}

	return nil
}

func JSONFormatter() log.Formatter {
	return &log.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		FieldMap: log.FieldMap{
			log.FieldKeyTime:  "time",
			log.FieldKeyLevel: "level",
			log.FieldKeyMsg:   "message",
			log.FieldKeyFunc:  "caller",
		},
	}
}

func FromContext(ctx context.Context) *log.Entry {
	entry := log.WithContext(ctx)
	if traceID := requestctx.TraceID(ctx); traceID != "" {
		entry = entry.WithField("trace_id", traceID)
	}
	if clientIP := requestctx.ClientIP(ctx); clientIP != "" {
		entry = entry.WithField("client_ip", clientIP)
	}
	if language := requestctx.Language(ctx); language != "" {
		entry = entry.WithField("language", language)
	}
	return entry
}

func RotatingWriter(filename string, size string) (io.Writer, error) {
	prepareLogFile(filename)
	bytes, err := units.ParseByteSize(size)
	if err != nil {
		return nil, err
	}
	maxSize := bytes / units.Megabyte
	if maxSize == 0 {
		return nil, fmt.Errorf("log size %q must be at least 1MB", size)
	}
	return &lumberjack.Logger{
		Filename: filename,
		MaxSize:  maxSize,
		Compress: true,
	}, nil
}

func prepareLogFile(filename string) {
	if filename == "" || filename == "stdout" || filename == "stderr" {
		return
	}

	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "create log dir: %v\n", err)
		return
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open log file: %v\n", err)
		return
	}
	if err := file.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "close log file: %v\n", err)
	}
	if err := os.Chmod(filename, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "chmod log file: %v\n", err)
	}
}
