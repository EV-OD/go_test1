package config

import (
	"strings"

	apperrors "myapp/internal/errors"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DatabaseConfig struct {
	// URL      string `mapstructure:"url"`
	USER     string `mapstructure: "user"`
	PASSWORD string `mapstructure: "password"`
	DB       string `mapstructure: "db"`
	PORT     int    `mapstructure: "port"`
	HOST     string `mapstructure: "host"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read the file
	if err := viper.ReadInConfig(); err != nil {
		return nil, apperrors.NewWithErr(apperrors.CodeConfigError, "failed to read config file", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, apperrors.NewWithErr(apperrors.CodeConfigError, "unable to decode into struct", err)
	}

	return &cfg, nil
}
