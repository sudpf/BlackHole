package config

import (
	"bytes"
	"fmt"
	"os"
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

func Load(file string) (*Config, error) {
	appName, err := executableName()
	if err != nil {
		return nil, err
	}

	configFile, err := resolveConfigFile(file)
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

func (c *Config) ListenHTTP() string {
	return c.App.ListenHttp
}

func (c *Config) RequestTimeout() time.Duration {
	return c.App.RequestTimeout
}

func (c *Config) ShutdownTimeout() time.Duration {
	return c.App.ShutdownTimeout
}

func resolveConfigFile(file string) (string, error) {
	if strings.TrimSpace(file) == "" {
		return "", fmt.Errorf("config file is required")
	}
	if filepath.IsAbs(file) {
		return filepath.Clean(file), nil
	}

	absPath, err := filepath.Abs(file)
	if err != nil {
		return "", fmt.Errorf("resolve config file %q: %w", file, err)
	}
	return absPath, nil
}

func executableName() (string, error) {
	appPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("get executable path: %w", err)
	}
	return filepath.Base(appPath), nil
}

func validate(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("voidengine config is required")
	}
	if strings.TrimSpace(cfg.App.ListenHttp) == "" {
		return fmt.Errorf("app.listen_http is required")
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
