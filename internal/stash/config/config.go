package config

import (
	"BlackHole/internal/runtime/configpath"
	"bytes"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"gopkg.in/yaml.v2"
)

type Config struct {
	App         appConfig      `yaml:"app" json:"app"`
	Log         logConfig      `yaml:"log" json:"log"`
	Metrics     metricsConfig  `yaml:"metrics" json:"metrics,optional"`
	Clusters    []*ClusterConf `yaml:"clusters" json:"clusters"`
	GracePeriod time.Duration  `yaml:"grace_period" json:"grace_period,default=10s"`
	appLogFile  string
}

type appConfig struct {
	Listen          string        `toml:"listen" yaml:"listen" json:"listen,default=http://127.0.0.1:80"`
	RequestTimeout  time.Duration `toml:"request_timeout" yaml:"request_timeout" json:"request_timeout,default=30s"`
	ShutdownTimeout time.Duration `toml:"shutdown_timeout" yaml:"shutdown_timeout" json:"shutdown_timeout,default=10s"`
}

type logConfig struct {
	Level string `toml:"level" yaml:"level" json:"level,default=info"`
	Size  string `toml:"size" yaml:"size" json:"size,default=256m"`
	Dir   string `toml:"dir" yaml:"dir" json:"dir,default=logs"`
}

type metricsConfig struct {
	Enabled bool   `yaml:"enabled" json:"enabled,default=false"`
	Listen  string `yaml:"listen" json:"listen,default=http://127.0.0.1:9102"`
	Path    string `yaml:"path" json:"path,default=/metrics"`
}

type ConditionConf struct {
	Key   string `yaml:"key" json:"key"`
	Value string `yaml:"value" json:"value"`
	Type  string `yaml:"type" json:"type,default=match,options=match|contains"`
	Op    string `yaml:"op" json:"op,default=and,options=and|or"`
}

type ElasticSearchConf struct {
	Hosts         []string `yaml:"hosts" json:"hosts"`
	Index         string   `yaml:"index" json:"index"`
	DocType       string   `yaml:"doc_type" json:"doc_type,default=doc"`
	TimeZone      string   `yaml:"time_zone" json:"time_zone,optional"`
	MaxChunkBytes int      `yaml:"max_chunk_bytes" json:"max_chunk_bytes,default=15728640"` // default 15M
	Compress      bool     `yaml:"compress" json:"compress,default=false"`
	Username      string   `yaml:"username" json:"username,optional"`
	Password      string   `yaml:"password" json:"password,optional"`
}

type SyslogAddrConf struct {
	Protocol string   `yaml:"protocol" json:"protocol,default=Udp"`
	Address  string   `yaml:"address" json:"address"`
	Port     int      `yaml:"port" json:"port"`
	Columns  []string `yaml:"columns" json:"columns,optional"`
}

type SyslogOutputConf struct {
	Conditions  [][]*ConditionConf `yaml:"conditions" json:"conditions,optional"`
	SyslogAddrs []*SyslogAddrConf  `yaml:"syslog_addrs" json:"syslog_addrs,optional"`
}

type FilterConf struct {
	Action     string          `yaml:"action" json:"action,options=drop|remove_field|transfer"`
	Conditions []ConditionConf `yaml:"conditions" json:"conditions,optional"`
	Fields     []string        `yaml:"fields" json:"fields,optional"`
	Field      string          `yaml:"field" json:"field,optional"`
	Target     string          `yaml:"target" json:"target,optional"`
}

type SyslogServiceConf struct {
	Protocol   string `yaml:"protocol" json:"protocol,default=Udp"`
	Ssl        string `yaml:"ssl,omitempty" json:"ssl,options=on|off,default=off"`
	Address    string `yaml:"address,omitempty" json:"address,optional"`
	Port       int    `yaml:"port,omitempty" json:"port,optional"`
	Processors int    `yaml:"processors" json:"processors,default=2"`
}

type KafkaConf struct {
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

type ClickHouseAuthConf struct {
	Database string `yaml:"database" json:"database"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}

type ClickHouseConf struct {
	Addr           []string           `yaml:"addr" json:"addr"`
	Auth           ClickHouseAuthConf `yaml:"auth" json:"auth"`
	Table          string             `yaml:"table" json:"table"`
	Columns        []string           `yaml:"columns" json:"columns,optional"`
	FillNoneColumn bool               `yaml:"fill_none_column" json:"fill_none_column,default=true"`
	Interval       int64              `yaml:"interval" json:"interval,default=15"`
	MaxChunkBytes  int                `yaml:"max_chunk_bytes" json:"max_chunk_bytes,default=15728640"` // default 15M
}

type InputConf struct {
	Syslogs []*SyslogServiceConf `yaml:"syslogs,omitempty" json:"syslogs,optional"`
	Kafka   *KafkaConf           `yaml:"kafka,omitempty" json:"kafka,optional"`
}

type OutputConf struct {
	ElasticSearch *ElasticSearchConf  `yaml:"elasticsearch,omitempty" json:"elasticsearch,optional"`
	Syslogs       []*SyslogOutputConf `yaml:"syslogs,omitempty" json:"syslogs,optional"`
	Clickhouse    *ClickHouseConf     `yaml:"clickhouse,omitempty" json:"clickhouse,optional"`
}

type ClusterConf struct {
	Input   *InputConf   `yaml:"input,omitempty" json:"input,optional"`
	Filters []FilterConf `yaml:"filters,omitempty" json:"filters,optional"`
	Output  *OutputConf  `yaml:"output,omitempty" json:"output,optional"`
}

func Load(file string) (*Config, error) {
	appName, err := configpath.ExecutableName()
	if err != nil {
		return nil, err
	}

	configFile, err := configpath.Resolve(file)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := conf.Load(configFile, cfg); err != nil {
		return nil, fmt.Errorf("load stash config: %w", err)
	}
	if err := validate(cfg); err != nil {
		return nil, err
	}

	cfg.appLogFile = filepath.Join(cfg.LogDir(), appName+".log")
	return cfg, nil
}

func (c *Config) String() string {
	buf := new(bytes.Buffer)
	if err := yaml.NewEncoder(buf).Encode(c); err != nil {
		return fmt.Sprintf("encode stash config: %v", err)
	}

	return buf.String()
}

func (c *Config) AppLogFile() string {
	return c.appLogFile
}

func (c *Config) LogLevel() string {
	return c.Log.Level
}

func (c *Config) LogSize() string {
	return c.Log.Size
}

func (c *Config) LogDir() string {
	return c.Log.Dir
}

func (c *Config) ShutdownTimeout() time.Duration {
	return c.App.ShutdownTimeout
}

func (c *Config) MetricsEnabled() bool {
	return c.Metrics.Enabled
}

func (c *Config) MetricsListenAddress() string {
	listenURL, _ := url.Parse(strings.TrimSpace(c.Metrics.Listen))
	return listenURL.Host
}

func (c *Config) MetricsPath() string {
	return c.Metrics.Path
}

func validate(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("stash config is required")
	}
	listenURL, err := url.Parse(strings.TrimSpace(cfg.App.Listen))
	if err != nil || listenURL.Scheme == "" || listenURL.Host == "" {
		return fmt.Errorf("app.listen must be a valid URL")
	}
	if listenURL.Scheme != "http" {
		return fmt.Errorf("app.listen only supports http scheme")
	}
	if cfg.App.RequestTimeout <= 0 {
		return fmt.Errorf("app.request_timeout must be positive")
	}
	if cfg.App.ShutdownTimeout <= 0 {
		return fmt.Errorf("app.shutdown_timeout must be positive")
	}
	if strings.TrimSpace(cfg.Log.Level) == "" {
		return fmt.Errorf("log.level is required")
	}
	if strings.TrimSpace(cfg.Log.Size) == "" {
		return fmt.Errorf("log.size is required")
	}
	if strings.TrimSpace(cfg.Log.Dir) == "" {
		return fmt.Errorf("log.dir is required")
	}
	if cfg.GracePeriod <= 0 {
		return fmt.Errorf("grace_period must be positive")
	}
	if cfg.Metrics.Enabled {
		metricsURL, err := url.Parse(strings.TrimSpace(cfg.Metrics.Listen))
		if err != nil || metricsURL.Scheme == "" || metricsURL.Host == "" {
			return fmt.Errorf("metrics.listen must be a valid URL")
		}
		if metricsURL.Scheme != "http" {
			return fmt.Errorf("metrics.listen only supports http scheme")
		}
		if !strings.HasPrefix(cfg.Metrics.Path, "/") {
			return fmt.Errorf("metrics.path must start with /")
		}
	}
	return nil
}
