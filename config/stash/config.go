package config

import (
	"bytes"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

type logConfig struct {
	Level  string `toml:"level"`
	Output string `toml:"output"`
}

type tomlConfig struct {
	Title string
	Log   logConfig
}

var GlobalConfig tomlConfig

func GetConfig() string {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(GlobalConfig); err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

func ParseConfig(file string) error {
	if _, err := toml.DecodeFile(file, &GlobalConfig); err != nil {
		return err
	}

	return nil
}
