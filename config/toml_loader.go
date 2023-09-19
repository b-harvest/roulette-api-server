package config

import (
	"fmt"
	"io/ioutil"

	toml "github.com/pelletier/go-toml"
)

var (
	DefaultConfigPath = "./config.toml"
)

// Config defines all necessary configuration parameters.
type tomlConfig struct {
	DBConf        TomlDBConf                  `toml:"config_db"`
	PriceConf     TomlPriceConf               `toml:"price"`
}

type TomlDBConf struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
	User string `toml:"user"`
	Password string `toml:"password"`
	DBName string `toml:"db_name"`
}
type TomlPriceConf struct {
	Url string `toml:"url"`
}

func Load(configPath string) (*tomlConfig, error) {
	var cfg tomlConfig
	if configPath == "" {
		return nil, fmt.Errorf("empty configuration path")
	}

	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %s", err)
	}
	err = toml.Unmarshal(configData, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config: %s", err)
	}

	return &cfg, nil
}
