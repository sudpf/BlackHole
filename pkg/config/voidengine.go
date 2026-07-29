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
	Title    string         `toml:"title" yaml:"title" json:"title,optional"`
	App      appConfig      `toml:"app" yaml:"app" json:"app"`
	Log      logConfig      `toml:"log" yaml:"log" json:"log"`
	Database DatabaseConfig `toml:"database" yaml:"database" json:"database"`
}

var (
	GlobalVoidEngineConfig VoidEngineConfig
	appVoidEngineName      string
	appVoidEngineBaseDir   string
	appVoidEngineLogFile   string
	apiVoidEngineLogFile   string
)

func (c *VoidEngineConfig) String() string {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(c); err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

func (c *VoidEngineConfig) AppLogFile() string {
	return appVoidEngineLogFile
}

func (c *VoidEngineConfig) ApiLogFile() string {
	return apiVoidEngineLogFile
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

func GetVoidEngineConfig() *VoidEngineConfig {
	return &GlobalVoidEngineConfig
}

func ParseVoidEngineConfig(file string) error {
	// 获取当前执行文件的路径
	appPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("get executable path: %w", err)
	}
	appVoidEngineName := filepath.Base(appPath)

	// 获取绝对路径
	absPath, err := filepath.Abs(appPath)
	if err != nil {
		return fmt.Errorf("get absolute executable path: %w", err)
	}

	// 获取目录路径
	appVoidEngineBaseDir := filepath.Dir(absPath)

	if !filepath.IsAbs(file) {
		file = appVoidEngineBaseDir + "/../conf/" + file
	}

	if err := conf.Load(file, &GlobalVoidEngineConfig); err != nil {
		return fmt.Errorf("load voidengine config: %w", err)
	}

	if !filepath.IsAbs(GlobalVoidEngineConfig.Log.Dir) {
		GlobalVoidEngineConfig.Log.Dir = appVoidEngineBaseDir + "/../" + GlobalVoidEngineConfig.Log.Dir
	}

	appVoidEngineLogFile = GlobalVoidEngineConfig.Log.Dir + "/" + appVoidEngineName + ".log"
	apiVoidEngineLogFile = GlobalVoidEngineConfig.Log.Dir + "/" + "api.log"

	return nil
}
