package logger

import (
	"BlackHole/pkg/requestctx"
	"BlackHole/pkg/units"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLog(level string, output string, size string) error {
	log.SetFormatter(JSONFormatter())
	log.SetReportCaller(true)

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
	workingDir, _ := os.Getwd()

	return &log.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		FieldMap: log.FieldMap{
			log.FieldKeyTime:  "time",
			log.FieldKeyLevel: "level",
			log.FieldKeyMsg:   "message",
		},
		CallerPrettyfier: func(frame *runtime.Frame) (string, string) {
			file := frame.File
			if workingDir != "" {
				if relativePath, err := filepath.Rel(workingDir, frame.File); err == nil {
					file = relativePath
				}
			}
			return "", fmt.Sprintf("%s:%d", file, frame.Line)
		},
	}
}

func FromContext(ctx context.Context) *log.Entry {
	entry := log.WithContext(ctx)
	if traceID := requestctx.TraceID(ctx); traceID != "" {
		entry = entry.WithField("trace_id", traceID)
	}
	return entry
}

func RotatingWriter(filename string, size string) (io.Writer, error) {
	if err := prepareLogFile(filename); err != nil {
		return nil, err
	}
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

func prepareLogFile(filename string) error {
	if filename == "" || filename == "stdout" || filename == "stderr" {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return fmt.Errorf("create log dir: %w", err)
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("open log file: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close log file: %w", err)
	}
	if err := os.Chmod(filename, 0644); err != nil {
		return fmt.Errorf("chmod log file: %w", err)
	}
	return nil
}
