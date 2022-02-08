// Package config - for working with yaml config file
package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"project/pkg/logging"
	"sync"
	"time"
)

// Config structure
type Config struct {
	IsDebug *bool `yaml:"is_debug" env-required:"true"`
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"localhost"`
		Port   string `yaml:"port" env-default:"8080"`
	} `yaml:"listen"`
	Timeout struct {
		Write time.Duration `yaml:"write" env-default:"15"`
		Read  time.Duration `yaml:"read" env-default:"15"`
	}
	MongoDB struct {
		Host       string `json:"host"`
		Port       string `json:"port"`
		Database   string `json:"database"`
		AuthDB     string `json:"auth_db"`
		Username   string `json:"username"`
		Password   string `json:"password"`
		Collection string `json:"collection"`
	} `json:"mongodb"`
}

var instance *Config
var once sync.Once

// GetConfig - Get data from config file
func GetConfig() *Config {
	once.Do(func() {
		// Logging
		logger := logging.GetLogger()
		logger.Info("read application configuration")
		instance = &Config{}

		// Get data from config with cleanenv
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
