package config

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"log"
)

func Load() (err error) {
	// oauth config
	if _, err := toml.DecodeFile("auth_config.toml", &ConfigOAuth); err != nil {
		return err
	}
	if marshal, err := json.Marshal(ConfigOAuth); err == nil {
		log.Println(string(marshal))
	}
	// web config
	if _, err := toml.DecodeFile("web_config.toml", &ConfigWeb); err != nil {
		return err
	}
	if marshal, err := json.Marshal(ConfigWeb); err == nil {
		log.Println(string(marshal))
	}
	return nil
}
