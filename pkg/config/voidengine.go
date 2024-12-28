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

func GetVoidEngineConfig() *VoidEngineConfig {
	return &GlobalVoidEngineConfig
}

func ParseVoidEngineConfig(file string) error {
	// 获取当前执行文件的路径
	appPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	appVoidEngineName := filepath.Base(appPath)

	// 获取绝对路径
	absPath, err := filepath.Abs(appPath)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	// 获取目录路径
	appVoidEngineBaseDir := filepath.Dir(absPath)

	if !filepath.IsAbs(file) {
		file = appVoidEngineBaseDir + "/../conf/" + file
	}

	if _, err := toml.DecodeFile(file, &GlobalVoidEngineConfig); err != nil {
		return err
	}

	if !filepath.IsAbs(GlobalVoidEngineConfig.Log.Dir) {
		GlobalVoidEngineConfig.Log.Dir = appVoidEngineBaseDir + "/../" + GlobalVoidEngineConfig.Log.Dir
	}

	appVoidEngineLogFile = GlobalVoidEngineConfig.Log.Dir + "/" + appVoidEngineName + ".log"
	apiVoidEngineLogFile = GlobalVoidEngineConfig.Log.Dir + "/" + "api.log"

	return nil
}
