package config

import "time"

type appConfig struct {
	ListenHttp      string        `toml:"listen_http" yaml:"listen_http" json:"listen_http,default=127.0.0.1:80"`
	ListenHttps     string        `toml:"listen_https" yaml:"listen_https" json:"listen_https,optional"`
	RequestTimeout  time.Duration `toml:"request_timeout" yaml:"request_timeout" json:"request_timeout,default=30s"`
	ShutdownTimeout time.Duration `toml:"shutdown_timeout" yaml:"shutdown_timeout" json:"shutdown_timeout,default=10s"`
}

type logConfig struct {
	Level string `toml:"level" yaml:"level" json:"level,default=info"`
	Size  string `toml:"size" yaml:"size" json:"size,default=256m"`
	Dir   string `toml:"dir" yaml:"dir" json:"dir,default=logs"`
	Gin   string `toml:"gin" yaml:"gin" json:"gin,optional"`
}

type MySQLConfig struct {
	Debug bool   `toml:"debug" yaml:"debug" json:"debug"`
	Log   string `toml:"log" yaml:"log" json:"log"`
	Link  string `toml:"link" yaml:"link" json:"link"`
}

type ClickHouseConfig struct {
	Debug bool   `toml:"debug" yaml:"debug" json:"debug"`
	Log   string `toml:"log" yaml:"log" json:"log"`
	Link  string `toml:"link" yaml:"link" json:"link"`
}

type DatabaseConfig struct {
	MySQL      *MySQLConfig      `toml:"mysql" yaml:"mysql" json:"mysql,optional"`
	ClickHouse *ClickHouseConfig `toml:"clickhouse" yaml:"clickhouse" json:"clickhouse,optional"`
}
