package db

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database interface {
	Connect(connectionString string) (*gorm.DB, error)
	Close() error
	CreateTable(model interface{}) error
	Query(model interface{}, conditions map[string]interface{}) (*gorm.DB, error)
	QueryEx(model interface{}, conditions interface{}) (*gorm.DB, error)
	Insert(model interface{}) error
	Update(model interface{}, conditions map[string]interface{}) error
	Delete(model interface{}, conditions map[string]interface{}) error
}

type logrusAdapter struct {
	logger *logrus.Logger
}

func NewLogrusAdapter(l *logrus.Logger) *logrusAdapter {
	return &logrusAdapter{logger: l}
}

func (l *logrusAdapter) LogMode(level logger.LogLevel) logger.Interface {
	// Implement logic to set log level if necessary
	return l
}

func (l *logrusAdapter) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger.WithContext(ctx).Infof(msg, data...)
}

func (l *logrusAdapter) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger.WithContext(ctx).Warnf(msg, data...)
}

func (l *logrusAdapter) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger.WithContext(ctx).Errorf(msg, data...)
}

func (l *logrusAdapter) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	entry := l.logger.WithContext(ctx).WithFields(logrus.Fields{
		"sql":     sql,
		"rows":    rows,
		"elapsed": elapsed.Seconds(),
	})

	if err != nil {
		entry.Errorf("trace error: %v", err)
	} else {
		entry.Infof("trace")
	}
}

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var logLevel string
	switch entry.Level {
	case logrus.InfoLevel:
		logLevel = "INFO"
	case logrus.WarnLevel:
		logLevel = "WARN"
	case logrus.ErrorLevel:
		logLevel = "ERROR"
	case logrus.DebugLevel:
		logLevel = "DEBUG"
	default:
		logLevel = "UNKNOWN"
	}

	// Format the log entry
	logMsg := fmt.Sprintf("%s [%s] %s | Elapsed: %v, Rows: %d, SQL: %s\n",
		entry.Time.Format(time.RFC3339),
		logLevel,
		entry.Message,
		entry.Data["elapsed"], // Assuming elapsed_ms is included in data
		entry.Data["rows"],
		entry.Data["sql"],
	)
	return []byte(logMsg), nil
}

func StructToConditions(v interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct but got %s", value.Kind())
	}

	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldName := typ.Field(i).Name

		// 跳过零值
		if field.IsZero() {
			continue
		}

		// 添加到结果 map
		result[fieldName] = field.Interface()
	}
	return result, nil
}
