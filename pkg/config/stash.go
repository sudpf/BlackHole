package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"gopkg.in/yaml.v2"
)

type (
	ConditionConf struct {
		Key   string `yaml:"key" json:"key"`
		Value string `yaml:"value" json:"value"`
		Type  string `yaml:"type" json:"type,default=match,options=match|contains"`
		Op    string `yaml:"op" json:"op,default=and,options=and|or"`
	}

	ElasticSearchConf struct {
		Hosts         []string `yaml:"hosts" json:"hosts"`
		Index         string   `yaml:"index" json:"index"`
		DocType       string   `yaml:"doc_type" json:"doc_type,default=doc"`
		TimeZone      string   `yaml:"time_zone" json:"time_zone,optional"`
		MaxChunkBytes int      `yaml:"max_chunk_bytes" json:"max_chunk_bytes,default=15728640"` // default 15M
		Compress      bool     `yaml:"compress" json:"compress,default=false"`
		Username      string   `yaml:"username" json:"username,optional"`
		Password      string   `yaml:"password" json:"password,optional"`
	}

	SyslogAddrConf struct {
		Protocol string   `yaml:"protocol" json:"protocol,default=Udp"`
		Address  string   `yaml:"address" json:"address"`
		Port     int      `yaml:"port" json:"port"`
		Columns  []string `yaml:"columns" json:"columns,optional"`
	}

	SyslogOutputConf struct {
		Conditions  [][]*ConditionConf `yaml:"conditions" json:"conditions,optional"`
		SyslogAddrs []*SyslogAddrConf  `yaml:"syslog_addrs" json:"syslog_addrs,optional"`
	}

	FilterConf struct {
		Action     string          `yaml:"action" json:"action,options=drop|remove_field|transfer"`
		Conditions []ConditionConf `yaml:"conditions" json:"conditions,optional"`
		Fields     []string        `yaml:"fields" json:"fields,optional"`
		Field      string          `yaml:"field" json:"field,optional"`
		Target     string          `yaml:"target" json:"target,optional"`
	}

	SyslogServiceConf struct {
		Protocol   string `yaml:"protocol" json:"protocol,default=Udp"`
		Ssl        string `yaml:"ssl,omitempty" json:"ssl,options=on|off,default=off"`
		Address    string `yaml:"address,omitempty" json:"address,optional"`
		Port       int    `yaml:"port,omitempty" json:"port,optional"`
		Processors int    `yaml:"processors" json:"processors,default=2"`
	}

	KafkaConf struct {
		service.ServiceConf
		Brokers    []string `yaml:"brokers" json:"brokers"`
		Group      string   `yaml:"group" json:"group"`
		Topics     []string `yaml:"topics" json:"topics"`
		Offset     string   `yaml:"offset" json:"offset,options=first|last,default=last"`
		Conns      int      `yaml:"conns" json:"conns,default=1"`
		Consumers  int      `yaml:"consumers" json:"consumers,default=8"`
		Processors int      `yaml:"processors" json:"processors,default=8"`
		MinBytes   int      `yaml:"min_bytes" json:"min_bytes,default=10240"`    // default 10K
		MaxBytes   int      `yaml:"max_bytes" json:"max_bytes,default=10485760"` // default 10M
		Username   string   `yaml:"username" json:"username,optional"`
		Password   string   `yaml:"password" json:"password,optional"`
	}

	ClickHouseAuthConf struct {
		Database string `yaml:"database" json:"database"`
		Username string `yaml:"username" json:"username"`
		Password string `yaml:"password" json:"password"`
	}

	ClickHouseConf struct {
		Addr           []string           `yaml:"addr" json:"addr"`
		Auth           ClickHouseAuthConf `yaml:"auth" json:"auth"`
		Table          string             `yaml:"table" json:"table"`
		Columns        []string           `yaml:"columns" json:"columns,optional"`
		FillNoneColumn bool               `yaml:"fill_none_column" json:"fill_none_column,default=true"`
		Interval       int64              `yaml:"interval" json:"interval,default=15"`
		MaxChunkBytes  int                `yaml:"max_chunk_bytes" json:"max_chunk_bytes,default=15728640"` // default 15M
	}

	InputConf struct {
		Syslogs []*SyslogServiceConf `yaml:"syslogs,omitempty" json:"syslogs,optional"`
		Kafka   *KafkaConf           `yaml:"kafka,omitempty" json:"kafka,optional"`
	}

	OutputConf struct {
		ElasticSearch *ElasticSearchConf  `yaml:"elasticsearch,omitempty" json:"elasticsearch,optional"`
		Syslogs       []*SyslogOutputConf `yaml:"syslogs,omitempty" json:"syslogs,optional"`
		Clickhouse    *ClickHouseConf     `yaml:"clickhouse,omitempty" json:"clickhouse,optional"`
	}

	ClusterConf struct {
		Input   *InputConf   `yaml:"input,omitempty" json:"input,optional"`
		Filters []FilterConf `yaml:"filters,omitempty" json:"filters,optional"`
		Output  *OutputConf  `yaml:"output,omitempty" json:"output,optional"`
	}

	StashConfig struct {
		App         appConfig      `yaml:"app" json:"app"`
		Log         logConfig      `yaml:"log" json:"log"`
		Clusters    []*ClusterConf `yaml:"clusters" json:"clusters"`
		GracePeriod time.Duration  `yaml:"grace_period" json:"grace_period,default=10s"`
		appLogFile  string
		apiLogFile  string
	}
)

func (c *StashConfig) String() string {
	buf := new(bytes.Buffer)
	if err := yaml.NewEncoder(buf).Encode(c); err != nil {
		return fmt.Sprintf("encode stash config: %v", err)
	}

	return buf.String()
}

func (c *StashConfig) AppLogFile() string {
	return c.appLogFile
}

func (c *StashConfig) ApiLogFile() string {
	return c.apiLogFile
}

func (c *StashConfig) LogLevel() string {
	return c.Log.Level
}

func (c *StashConfig) LogSize() string {
	return c.Log.Size
}

func (c *StashConfig) LogDir() string {
	return c.Log.Dir
}

func LoadStashConfig(file string) (*StashConfig, error) {
	appPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("get executable path: %w", err)
	}
	appName := filepath.Base(appPath)

	absPath, err := filepath.Abs(appPath)
	if err != nil {
		return nil, fmt.Errorf("get absolute executable path: %w", err)
	}
	appBaseDir := filepath.Dir(absPath)

	if !filepath.IsAbs(file) {
		file = filepath.Join(appBaseDir, "..", "conf", file)
	}

	cfg := &StashConfig{}
	if err := conf.Load(file, cfg); err != nil {
		return nil, fmt.Errorf("load stash config: %w", err)
	}

	if !filepath.IsAbs(cfg.Log.Dir) {
		cfg.Log.Dir = filepath.Join(appBaseDir, "..", cfg.Log.Dir)
	}

	cfg.appLogFile = filepath.Join(cfg.Log.Dir, appName+".log")
	cfg.apiLogFile = filepath.Join(cfg.Log.Dir, appName+"api.log")

	return cfg, nil
}
