package config

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

type VoidEngineConfig struct {
	Title    string
	App      appConfig
	Log      logConfig
	Database DatabaseConfig `toml:"database"`
}

var (
	GlobalVoidEngineConfig VoidEngineConfig
	appName                string
	appBaseDir             string
	appLogFile             string
	apiLogFile             string
)

func (c *VoidEngineConfig) String() string {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(c); err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

func (c *VoidEngineConfig) AppLogFile() string {
	return appLogFile
}

func (c *VoidEngineConfig) ApiLogFile() string {
	return apiLogFile
}

func (c *VoidEngineConfig) LogLevel() string {
	return c.Log.Level
}

func (c *VoidEngineConfig) LogDir() string {
	return c.Log.Dir
}

func GetVoidEngineConfig() *VoidEngineConfig {
	return &GlobalVoidEngineConfig
}

func ParseVoidEngineConfig(file string) error {
	// 获取当前执行文件的路径
	appPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	appName := filepath.Base(appPath)

	// 获取绝对路径
	absPath, err := filepath.Abs(appPath)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	// 获取目录路径
	appBaseDir := filepath.Dir(absPath)

	if !filepath.IsAbs(file) {
		file = appBaseDir + "/../conf/" + file
	}

	if _, err := toml.DecodeFile(file, &GlobalVoidEngineConfig); err != nil {
		return err
	}

	if !filepath.IsAbs(GlobalVoidEngineConfig.Log.Dir) {
		GlobalVoidEngineConfig.Log.Dir = appBaseDir + "/../" + GlobalVoidEngineConfig.Log.Dir
	}

	appLogFile = GlobalVoidEngineConfig.Log.Dir + "/" + appName + ".log"
	apiLogFile = GlobalVoidEngineConfig.Log.Dir + "/" + "api.log"

	return nil
}
