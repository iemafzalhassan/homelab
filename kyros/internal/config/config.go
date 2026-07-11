package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port     string
	LogLevel string
	Env      string
}

func Load(serviceName string) *Config {
	v := viper.New()
	v.SetEnvPrefix("KYROS")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("port", "8080")
	v.SetDefault("log_level", "info")
	v.SetDefault("env", "development")

	return &Config{
		Port:     v.GetString("port"),
		LogLevel: v.GetString("log_level"),
		Env:      v.GetString("env"),
	}
}
