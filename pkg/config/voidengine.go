package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"github.com/zeromicro/go-zero/core/conf"
)

type VoidEngineConfig struct {
	Title      string         `toml:"title" yaml:"title" json:"title,optional"`
	App        appConfig      `toml:"app" yaml:"app" json:"app"`
	Log        logConfig      `toml:"log" yaml:"log" json:"log"`
	Database   DatabaseConfig `toml:"database" yaml:"database" json:"database"`
	appLogFile string
	apiLogFile string
}

func (c *VoidEngineConfig) String() string {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(c); err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

func (c *VoidEngineConfig) AppLogFile() string {
	return c.appLogFile
}

func (c *VoidEngineConfig) ApiLogFile() string {
	return c.apiLogFile
}

func (c *VoidEngineConfig) LogLevel() string {
	return c.Log.Level
}

func (c *VoidEngineConfig) LogDir() string {
	return c.Log.Dir
}

func (c *VoidEngineConfig) ListenHTTP() string {
	return c.App.ListenHttp
}

func (c *VoidEngineConfig) ShutdownTimeout() time.Duration {
	return c.App.ShutdownTimeout
}

func LoadVoidEngineConfig(file string) (*VoidEngineConfig, error) {
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

	cfg := &VoidEngineConfig{}
	if err := conf.Load(file, cfg); err != nil {
		return nil, fmt.Errorf("load voidengine config: %w", err)
	}

	if !filepath.IsAbs(cfg.Log.Dir) {
		cfg.Log.Dir = filepath.Join(appBaseDir, "..", cfg.Log.Dir)
	}

	cfg.appLogFile = filepath.Join(cfg.Log.Dir, appName+".log")
	cfg.apiLogFile = filepath.Join(cfg.Log.Dir, "api.log")

	return cfg, nil
}
