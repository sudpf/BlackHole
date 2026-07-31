package db

import (
	"BlackHole/pkg/requestctx"
	"context"
	"encoding/json"
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
	CreateTable(model ...interface{}) error
	Query(ctx context.Context, model interface{}, conditions map[string]interface{}) (*gorm.DB, error)
	QueryEx(ctx context.Context, model interface{}, conditions interface{}) (*gorm.DB, error)
	Insert(ctx context.Context, model interface{}) error
	Update(ctx context.Context, model interface{}, conditions map[string]interface{}) error
	Delete(ctx context.Context, model interface{}, conditions map[string]interface{}) error
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
	l.entry(ctx).Infof(msg, data...)
}

func (l *logrusAdapter) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.entry(ctx).Warnf(msg, data...)
}

func (l *logrusAdapter) Error(ctx context.Context, msg string, data ...interface{}) {
	l.entry(ctx).Errorf(msg, data...)
}

func (l *logrusAdapter) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	entry := l.entry(ctx).WithFields(logrus.Fields{
		"sql":        sql,
		"rows":       rows,
		"elapsed_ms": float64(elapsed.Microseconds()) / 1000,
	})

	if err != nil {
		entry.WithError(err).Error("trace")
	} else {
		entry.Infof("trace")
	}
}

func (l *logrusAdapter) entry(ctx context.Context) *logrus.Entry {
	entry := l.logger.WithContext(ctx)
	if traceID := requestctx.TraceID(ctx); traceID != "" {
		entry = entry.WithField("trace_id", traceID)
	} else {
		entry = entry.WithField("trace_id", "system")
	}
	return entry
}

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	fields := logrus.Fields{
		"time":       entry.Time.Format("2006-01-02T15:04:05.000Z07:00"),
		"level":      entry.Level.String(),
		"trace_id":   entry.Data["trace_id"],
		"elapsed_ms": entry.Data["elapsed_ms"],
		"rows":       entry.Data["rows"],
		"sql":        entry.Data["sql"],
	}
	if err, ok := entry.Data[logrus.ErrorKey].(error); ok && err != nil {
		fields["error"] = err.Error()
	}
	if entry.Message != "" && entry.Message != "trace" {
		fields["message"] = entry.Message
	}

	message, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}
	return append(message, '\n'), nil
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
