package output

import (
	"BlackHole/pkg/config"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	log "github.com/sirupsen/logrus"
	"github.com/zeromicro/go-zero/core/executors"
)

type (
	ChWriter struct {
		client         clickhouse.Conn
		ctx            context.Context
		columns        []string
		columnsType    map[string]string
		query          string
		database       string
		table          string
		inserter       *executors.ChunkExecutor
		fillNoneColumn bool
	}

	ValueWithIndex struct {
		val []interface{}
	}

	ClickHouseDescType struct {
		Name              string `db:"name"`
		Type              string `db:"type"`
		DefaultType       string `db:"default_type"`
		DefaultExpression string `db:"default_expression"`
		Comment           string `db:"comment"`
		CodecExpression   string `db:"codec_expression"`
		TTLExpression     string `db:"ttl_expression"`
	}
)

func NewClickHouseWriter(c *config.ClickHouseConf) (*ChWriter, error) {
	client, err := clickhouse.Open(&clickhouse.Options{
		Addr: c.Addr,
		Auth: clickhouse.Auth{
			Database: c.Auth.Database,
			Username: c.Auth.Username,
			Password: c.Auth.Password,
		},
	})
	if err != nil {
		log.Warnf("open clickhouse err: %s", err)
		return nil, err
	}
	var query string
	if len(c.Columns) > 0 {
		query = "INSERT INTO " + c.Auth.Database + "." + c.Table + " (" + strings.Join(c.Columns, ", ") + ")"
	} else {
		query = "INSERT INTO " + c.Auth.Database + "." + c.Table
	}

	writer := ChWriter{
		client:         client,
		ctx:            context.Background(),
		database:       c.Auth.Database,
		table:          c.Table,
		query:          query,
		columns:        c.Columns,
		fillNoneColumn: c.FillNoneColumn,
	}
	writer.inserter = executors.NewChunkExecutor(writer.execute, executors.WithChunkBytes(c.MaxChunkBytes), executors.WithFlushInterval(time.Duration(c.Interval)*time.Second))

	for i := 0; i < 30; i++ {
		if err := writer.clickhouseColumns(); err == nil {
			return &writer, nil
		}
		log.Warnf("Init clickhouseColumns error: %v", err)
		time.Sleep(time.Second * 5)
	}
	return nil, fmt.Errorf("Init clickhouseColumns error")
}

func (w *ChWriter) clickhouseTableDesc() ([]*ClickHouseDescType, error) {
	var descTypes []*ClickHouseDescType

	rows, err := w.client.Query(w.ctx, "DESC "+w.database+"."+w.table)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var descType ClickHouseDescType
		err = rows.Scan(&descType.Name, &descType.Type, &descType.DefaultType, &descType.DefaultExpression, &descType.Comment, &descType.CodecExpression, &descType.TTLExpression)
		if err != nil {
			return nil, err
		}
		descTypes = append(descTypes, &descType)
	}
	defer rows.Close()

	return descTypes, nil
}

func (w *ChWriter) clickhouseColumns() error {
	descTypes, err := w.clickhouseTableDesc()
	if err != nil {
		return err
	}

	columnType := make(map[string]string)
	for _, item := range descTypes {
		columnType[item.Name] = item.Type
	}
	w.columnsType = columnType

	if len(w.columns) == 0 {
		for _, desc := range descTypes {
			w.columns = append(w.columns, desc.Name)
		}
	}

	for _, column := range w.columns {
		_, ok := columnType[column]
		if !ok {
			return errors.New(column + " not in " + w.database + " " + w.table)
		}
	}

	log.Debugf("clickhouse colum type: %v", w.columnsType)
	return nil
}

func (w *ChWriter) Write(val map[string]interface{}) error {
	if w == nil {
		return errors.New("invalid writer")
	}

	v, err := w.PrepareData(val)
	if err != nil {
		log.Warnf("PrepareData error:%v", err)
		return nil
	}

	return w.inserter.Add(ValueWithIndex{
		val: v,
	}, len(val))
}

func (w *ChWriter) PrepareData(val map[string]interface{}) ([]interface{}, error) {
	result := make([]interface{}, len(w.columns))
	for index, column := range w.columns {
		c, ok := val[column]
		if !ok {
			if w.fillNoneColumn {
				var value interface{}
				c = value
			} else {
				return nil, errors.New(column + " not in data")
			}
		}
		v, err := toClickhouseType(c, w.columnsType[column])
		if err != nil {
			log.Warningf("Clickhouse [%v] Type error", c)
			return nil, err
		}
		result[index] = v
	}

	log.Debugf("prepare data:%v", result)
	return result, nil
}

func (w *ChWriter) execute(vals []interface{}) {
	bulk, err := w.client.PrepareBatch(w.ctx, w.query)
	if err != nil {
		log.Warnf("PrepareBatch err: %v", err)
		return
	}

	for _, val := range vals {
		err = bulk.Append(val.(ValueWithIndex).val...)
		if err != nil {
			log.Warnf("Bulk append err: %v", err)
			return
		}
	}

	err = bulk.Send()
	if err != nil {
		log.Warnf("Bulk send err: %v", err)
		return
	}

	log.Debugf("Bulk send data success")
}
