package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name string
		Port int
	}
	Database struct {
		Host     string
		User     string
		Password string
	}
	Cache struct {
		File    string
		Storage string
	}
}

func LoadConfig() (*Config, error) {
	viper.Reset()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../config")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("whatsapi")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func ClearViperCache() {
	viper.Reset()
}
