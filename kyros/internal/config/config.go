package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config holds runtime configuration for a Kyros Go service. Values come from
// environment variables prefixed with KYROS_ (e.g. KYROS_DATABASE_URL).
type Config struct {
	Port               string
	LogLevel           string
	Env                string
	RegistryBackendURL string
	AuthRealm          string
	// KeycloakIssuer is the full issuer URL of the Keycloak realm, e.g.
	//   https://sso.iemafzalhassan.tech/realms/kyros
	// Both the API service and the registry proxy need this to validate
	// bearer tokens. They MUST agree — if they don't, the dashboard token
	// won't be accepted by the registry, and vice versa.
	KeycloakIssuer string
	DatabaseURL    string
}

func Load(serviceName string) *Config {
	v := viper.New()
	v.SetEnvPrefix("KYROS")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("port", "8080")
	v.SetDefault("log_level", "info")
	v.SetDefault("env", "development")

	v.SetDefault("registry_backend_url", "http://registry:5000")
	v.SetDefault("auth_realm", "Kyros Registry")
	v.SetDefault("database_url", "postgres://kyros:kyros@localhost:5432/kyros?sslmode=disable")

	// KeycloakIssuer intentionally has no default. The api service falls back
	// to a local dev issuer (see cmd/api/main.go) but the registry service
	// requires an explicit issuer — see cmd/registry/main.go and the
	// comments in internal/registry/auth.go about fail-closed behaviour.
	issuer := strings.TrimSpace(v.GetString("keycloak_issuer"))
	if issuer == "" {
		issuer = strings.TrimSpace(os.Getenv("KEYCLOAK_ISSUER"))
	}

	return &Config{
		Port:               v.GetString("port"),
		LogLevel:           v.GetString("log_level"),
		Env:                v.GetString("env"),
		RegistryBackendURL: v.GetString("registry_backend_url"),
		AuthRealm:          v.GetString("auth_realm"),
		KeycloakIssuer:     issuer,
		DatabaseURL:        v.GetString("database_url"),
	}
}
