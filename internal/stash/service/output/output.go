package output

import (
	"BlackHole/pkg/config"
	"errors"
	"fmt"
)

type (
	Writer interface {
		Write(val map[string]interface{}) error
		Close() error
	}
)

func NewWriters(o *config.OutputConf) ([]Writer, error) {
	var writers []Writer
	closeWriters := func() {
		_ = CloseWriters(writers)
	}

	if o == nil {
		return writers, nil
	}

	if o.ElasticSearch != nil {
		writer, err := NewElasticSearchWriter(o.ElasticSearch)
		if err != nil {
			closeWriters()
			return nil, fmt.Errorf("initialize elasticsearch writer: %w", err)
		}
		writers = append(writers, writer)
	}

	if o.Clickhouse != nil && len(o.Clickhouse.Addr) > 0 {
		writer, err := NewClickHouseWriter(o.Clickhouse)
		if err != nil {
			closeWriters()
			return nil, fmt.Errorf("initialize clickhouse writer: %w", err)
		}
		writers = append(writers, writer)
	}

	if len(o.Syslogs) > 0 {
		for i, syslog := range o.Syslogs {
			writer, err := NewSyslogWriter(syslog)
			if err != nil {
				closeWriters()
				return nil, fmt.Errorf("initialize syslog writer %d: %w", i, err)
			}
			writers = append(writers, writer)
		}
	}

	return writers, nil
}

func CloseWriters(writers []Writer) error {
	var err error
	for _, writer := range writers {
		err = errors.Join(err, writer.Close())
	}

	return err
}
