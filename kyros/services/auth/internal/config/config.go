package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config holds runtime configuration for the auth service.
type Config struct {
	Port               string
	LogLevel           string
	Env                string
	DatabaseURL        string
	JWTSecret          string
	AccessTokenExpiry  string
	RefreshTokenExpiry string
	BCryptCost         int
	ServiceName        string
	OTELCollectorURL   string
}

// Load reads configuration from environment variables and returns a Config struct.
func Load() *Config {
	v := viper.New()
	v.SetEnvPrefix("KYROS_AUTH")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Set defaults
	v.SetDefault("port", "8080")
	v.SetDefault("log_level", "info")
	v.SetDefault("env", "development")
	v.SetDefault("database_url", "postgres://kyros:kyros@localhost:5432/kyros?sslmode=disable")
	v.SetDefault("jwt_secret", "please-change-me-to-a-secure-secret")
	v.SetDefault("access_token_expiry", "15m")
	v.SetDefault("refresh_token_expiry", "720h") // 30 days
	v.SetDefault("bcrypt_cost", 12)
	v.SetDefault("service_name", "auth")
	v.SetDefault("otel_collector_url", "")

	return &Config{
		Port:               v.GetString("port"),
		LogLevel:           v.GetString("log_level"),
		Env:                v.GetString("env"),
		DatabaseURL:        v.GetString("database_url"),
		JWTSecret:          v.GetString("jwt_secret"),
		AccessTokenExpiry:  v.GetString("access_token_expiry"),
		RefreshTokenExpiry: v.GetString("refresh_token_expiry"),
		BCryptCost:         v.GetInt("bcrypt_cost"),
		ServiceName:        v.GetString("service_name"),
		OTELCollectorURL:   v.GetString("otel_collector_url"),
	}
}