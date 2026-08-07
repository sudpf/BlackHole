package config

import (
	"BlackHole/internal/runtime/configpath"
	"bytes"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/zeromicro/go-zero/core/conf"
)

type Config struct {
	Title      string         `toml:"title" yaml:"title" json:"title,optional"`
	App        appConfig      `toml:"app" yaml:"app" json:"app"`
	Log        logConfig      `toml:"log" yaml:"log" json:"log"`
	Database   DatabaseConfig `toml:"database" yaml:"database" json:"database"`
	appLogFile string
	apiLogFile string
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
		return nil, fmt.Errorf("load voidengine config: %w", err)
	}
	if err := validate(cfg); err != nil {
		return nil, err
	}

	return &Config{
		Title:      cfg.Title,
		App:        cfg.App,
		Log:        cfg.Log,
		Database:   cfg.Database,
		appLogFile: filepath.Join(cfg.LogDir(), appName+".log"),
		apiLogFile: filepath.Join(cfg.LogDir(), "api.log"),
	}, nil
}

func (c *Config) String() string {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(c); err != nil {
		return fmt.Sprintf("encode voidengine config: %v", err)
	}

	return buf.String()
}

func (c *Config) AppLogFile() string {
	return c.appLogFile
}

func (c *Config) ApiLogFile() string {
	return c.apiLogFile
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

func (c *Config) ListenAddress() string {
	listenURL, _ := url.Parse(strings.TrimSpace(c.App.Listen))
	return listenURL.Host
}

func (c *Config) RequestTimeout() time.Duration {
	return c.App.RequestTimeout
}

func (c *Config) ShutdownTimeout() time.Duration {
	return c.App.ShutdownTimeout
}

func validate(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("voidengine config is required")
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
	return nil
}
