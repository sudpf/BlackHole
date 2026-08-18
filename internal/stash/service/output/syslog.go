package output

import (
	"BlackHole/internal/stash/config"
	"BlackHole/internal/stash/service/filter"
	"context"
	"errors"
	"fmt"
	"log/syslog"
	"strings"

	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
)

type (
	SyslogWriteConf struct {
		Protocol string
		Address  string
		Port     int
		Columns  []string
		logger   *syslog.Writer
	}

	SyslogWriter struct {
		Filters    []filter.FilterFunc
		WriteConfs []*SyslogWriteConf
	}
)

func NewSyslogWriter(ctx context.Context, c *config.SyslogOutputConf) (*SyslogWriter, error) {
	w := &SyslogWriter{}
	if c == nil {
		return w, nil
	}
	w.Filters = createSyslogFilters(c.Conditions)
	closeWriter := func() {
		_ = w.Close(ctx)
	}

	for i, cSyslogAddr := range c.SyslogAddrs {
		if cSyslogAddr == nil {
			closeWriter()
			return nil, fmt.Errorf("syslog address %d is nil", i)
		}

		protocol := normalizeSyslogNetwork(cSyslogAddr.Protocol)
		address := syslogDialAddress(protocol, cSyslogAddr.Address, cSyslogAddr.Port)
		logger, err := syslog.Dial(protocol, address, syslog.LOG_INFO|syslog.LOG_LOCAL0, "stash")
		if err != nil {
			closeWriter()
			return nil, fmt.Errorf("dial syslog %s %s: %w", protocol, address, err)
		}

		w.WriteConfs = append(w.WriteConfs, &SyslogWriteConf{
			Protocol: protocol,
			Address:  cSyslogAddr.Address,
			Port:     cSyslogAddr.Port,
			Columns:  cSyslogAddr.Columns,
			logger:   logger,
		})
	}

	return w, nil
}

func createSyslogFilters(conditions [][]*config.ConditionConf) []filter.FilterFunc {
	filters := make([]filter.FilterFunc, 0, len(conditions))
	for _, group := range conditions {
		conds := make([]config.ConditionConf, 0, len(group))
		for _, cond := range group {
			if cond == nil {
				continue
			}
			conds = append(conds, *cond)
		}
		if len(conds) == 0 {
			continue
		}
		filters = append(filters, filter.DropFilter(conds))
	}
	return filters
}

func (w *SyslogWriter) PrepareData(columns []string, val map[string]interface{}) ([]interface{}, error) {
	result := make([]interface{}, len(columns))
	for index, column := range columns {
		v, ok := val[column]
		if !ok {
			var value interface{}
			v = value
		}

		result[index] = v
	}

	return result, nil
}

func (w *SyslogWriter) Write(ctx context.Context, val map[string]interface{}) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	filterMatch := false
	if len(w.Filters) > 0 {
		for _, filter := range w.Filters {
			if m := filter(val); m == nil {
				log.Debugf("Syslog [%v] conditions match", filter)
				filterMatch = true
				break
			} else {
				log.Debugf("Syslog miss match filter %v", filter)
			}
		}
	} else {
		log.Debugf("Syslog writer, No filters")
		filterMatch = true
	}

	if filterMatch == false {
		log.Debugf("Syslog writer, Miss match filter")
		return nil
	}

	var writeErr error
	for _, wc := range w.WriteConfs {
		if len(wc.Columns) == 0 {
			continue
		}

		prepared, err := w.PrepareData(wc.Columns, val)
		if err != nil {
			log.Warnf("PrepareData error:%v", err)
			return err
		}

		body := make(map[string]interface{}, len(wc.Columns))
		for i, v := range prepared {
			body[wc.Columns[i]] = v
		}

		dataString, err := jsoniter.MarshalToString(body)
		if err != nil {
			log.Warnf("Marshal syslog output error:%v", err)
			return err
		}

		err = wc.logger.Info(dataString)
		if err != nil {
			log.Warnf("Send syslog message error:%v", err)
			writeErr = errors.Join(writeErr, err)
		}
	}

	return writeErr
}

func (w *SyslogWriter) Close(ctx context.Context) error {
	var err error
	for _, writeConf := range w.WriteConfs {
		if writeConf.logger != nil {
			err = errors.Join(err, writeConf.logger.Close())
		}
	}

	return err
}

func normalizeSyslogNetwork(protocol string) string {
	switch strings.ToLower(protocol) {
	case "", "udp":
		return "udp"
	case "tcp":
		return "tcp"
	case "unix", "unixgram", "unixpacket":
		return strings.ToLower(protocol)
	default:
		return strings.ToLower(protocol)
	}
}

func syslogDialAddress(protocol, address string, port int) string {
	if strings.HasPrefix(protocol, "unix") {
		return address
	}

	return fmt.Sprintf("%s:%d", address, port)
}
