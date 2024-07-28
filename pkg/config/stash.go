package config

import (
	"bytes"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

type stashConfig struct {
	Title string
	App   appConfig
	Log   logConfig
}

var GlobalStashConfig stashConfig

func (c *stashConfig) String() string {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(c); err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

func ParseStashConfig(file string) error {
	if _, err := toml.DecodeFile(file, &GlobalStashConfig); err != nil {
		return err
	}

	return nil
}
