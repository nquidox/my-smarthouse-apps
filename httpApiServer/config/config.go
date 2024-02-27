package config

import (
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Hostname string `mapstructure:"hostname"`
	Port     int32  `mapstructure:"port"`
}

type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	DbPath string `mapstructure:"dbPath"`
}

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(path)

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
