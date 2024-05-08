package config

import (
	"be-capstone-project/src/internal/core/common_configs"
	"be-capstone-project/src/internal/core/logger"
	"encoding/json"
	"github.com/kelseyhightower/envconfig"
	"io/ioutil"
)

// Type is the type for config of
type Config struct {
	App   common_configs.App
	Store common_configs.Store
	Redis common_configs.Redis
	Kafka common_configs.KafkaSaramaConfig
}

// mustLoadConfig load config from env and file
func mustLoadConfig(confFile string, config *Config) {
	if confFile != "" {
		cont, err := ioutil.ReadFile(confFile)
		logger.FatalIfError(err)
		err = json.Unmarshal(cont, config)
		logger.FatalIfError(err)
	} else {
		err := envconfig.Process("", config)
		if err != nil {
			logger.FatalIfError(err)
		}
	}
}

// NewAppConfigs Init appconfig file. Can parse specific path to .yaml config file to load config
func NewAppConfigs(confFile string) (*Config, error) {
	var out Config
	mustLoadConfig(confFile, &out)
	return &out, nil
}
