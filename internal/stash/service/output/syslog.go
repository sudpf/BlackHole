package output

import (
	"BlackHole/internal/stash/service/filter"
	"BlackHole/pkg/config"
	"fmt"
	"log/syslog"

	log "github.com/sirupsen/logrus"
)

type (
	SyslogWriteConf struct {
		Protocol string
		Address  string
		Port     int
		Columns  []string
	}

	SyslogWriter struct {
		Filters    []filter.FilterFunc
		WriteConfs []*SyslogWriteConf
	}
)

func NewSyslogWriter(c *config.SyslogOutputConf) (*SyslogWriter, error) {
	w := &SyslogWriter{}
	for _, cSyslogAddr := range c.SyslogAddrs {
		w.WriteConfs = append(w.WriteConfs, &SyslogWriteConf{
			Protocol: cSyslogAddr.Protocol,
			Address:  cSyslogAddr.Address,
			Port:     cSyslogAddr.Port,
			Columns:  cSyslogAddr.Columns,
		})
	}

	return w, nil
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

func (w *SyslogWriter) Write(val map[string]interface{}) error {
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

	for _, wc := range w.WriteConfs {
		if len(wc.Columns) == 0 {
			continue
		}

		_val, err := w.PrepareData(wc.Columns, val)
		if err != nil {
			log.Warnf("PrepareData error:%v", err)
			return nil
		}

		var body string
		for i, v := range _val {
			if len(body) == 0 {
				body = fmt.Sprintf("\"%v\":\"%v\"", wc.Columns[i], v)
			} else {
				body = body + "," + fmt.Sprintf("\"%v\":\"%v\"", wc.Columns[i], v)
			}
		}
		dataString := fmt.Sprintf("{%s}", body)

		go func(address string, protocol string, port int, data string) {
			// 创建一个 Syslog logger，使用 UDP 协议连接到本地 Syslog 服务器
			logger, err := syslog.Dial(protocol, fmt.Sprintf("%s:%d", address, port), syslog.LOG_INFO|syslog.LOG_LOCAL0, "stash")
			if err != nil {
				log.Fatalf("Failed to dial syslog: %v", err)
			}
			defer logger.Close()

			// 发送一条 Syslog 消息
			err = logger.Info(data)
			if err != nil {
				log.Fatalf("Failed to send syslog message: %v", err)
			}
		}(wc.Address, wc.Protocol, wc.Port, dataString)
	}

	return nil
}
