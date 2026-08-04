package output

import (
	"BlackHole/internal/stash/config"
	"context"
	"errors"
	"fmt"
)

type (
	Writer interface {
		Write(ctx context.Context, val map[string]interface{}) error
		Close(ctx context.Context) error
	}
)

func NewWriters(ctx context.Context, o *config.OutputConf) ([]Writer, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}
	var writers []Writer
	closeWriters := func() {
		_ = CloseWriters(ctx, writers)
	}

	if o == nil {
		return writers, nil
	}

	if o.ElasticSearch != nil {
		writer, err := NewElasticSearchWriter(ctx, o.ElasticSearch)
		if err != nil {
			closeWriters()
			return nil, fmt.Errorf("initialize elasticsearch writer: %w", err)
		}
		writers = append(writers, writer)
	}

	if o.Clickhouse != nil && len(o.Clickhouse.Addr) > 0 {
		writer, err := NewClickHouseWriter(ctx, o.Clickhouse)
		if err != nil {
			closeWriters()
			return nil, fmt.Errorf("initialize clickhouse writer: %w", err)
		}
		writers = append(writers, writer)
	}

	if len(o.Syslogs) > 0 {
		for i, syslog := range o.Syslogs {
			writer, err := NewSyslogWriter(ctx, syslog)
			if err != nil {
				closeWriters()
				return nil, fmt.Errorf("initialize syslog writer %d: %w", i, err)
			}
			writers = append(writers, writer)
		}
	}

	return writers, nil
}

func CloseWriters(ctx context.Context, writers []Writer) error {
	var err error
	for _, writer := range writers {
		err = errors.Join(err, writer.Close(ctx))
	}

	return err
}
