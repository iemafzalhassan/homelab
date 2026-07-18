package main

import (
	"fmt"
	"os"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/kyros-platform/kyros/internal/auth"
	"github.com/kyros-platform/kyros/internal/config"
	"github.com/kyros-platform/kyros/internal/httpserver"
	"github.com/kyros-platform/kyros/internal/logger"
	"github.com/kyros-platform/kyros/internal/registry"
)

func main() {
	cfg := config.Load("registry")
	log := logger.New(cfg.LogLevel, cfg.Env)
	defer func() { _ = log.Sync() }()

	// The registry proxy is the auth layer for the entire OCI registry —
	// the backend cncf/distribution runs with auth disabled and trusts us.
	// A misconfigured auth path here means an open registry. Fail closed
	// if the Keycloak issuer is not configured.
	if cfg.KeycloakIssuer == "" {
		log.Fatal("KYROS_KEYCLOAK_ISSUER (or KEYCLOAK_ISSUER) is required for the registry service")
	}
	if cfg.AuthRealm == "" {
		log.Fatal("KYROS_AUTH_REALM is required for the registry service")
	}

	srv := httpserver.New(cfg.Port, log, "registry")

	validator := auth.NewValidator(cfg.KeycloakIssuer)

	proxy, err := registry.NewProxy(log.Sugar(), cfg.RegistryBackendURL)
	if err != nil {
		log.Fatal("failed to create proxy", zap.Error(err))
	}

	authMW := registry.NewAuthMiddleware(log.Sugar(), cfg.AuthRealm, validator)

	// /v2/* is the OCI Distribution API. Every endpoint here is guarded.
	srv.Router.Route("/v2", func(r chi.Router) {
		r.Use(authMW.RequireAuth)
		// Send everything to the proxy, stripping /v2 is NOT needed as the backend expects /v2/
		r.HandleFunc("/*", proxy.ServeHTTP)
		r.HandleFunc("/", proxy.ServeHTTP)
	})

	// Build-time sanity check: log which Keycloak realm we're validating against.
	log.Info("registry proxy ready",
		zap.String("issuer", cfg.KeycloakIssuer),
		zap.String("realm", cfg.AuthRealm),
		zap.String("backend", cfg.RegistryBackendURL),
	)

	// Helpful fatal if the user misconfigured env. Otherwise the validator
	// would fail at first request with a less obvious error.
	if v := os.Getenv("KYROS_AUTH_SECRET"); v != "" {
		fmt.Fprintln(os.Stderr,
			"WARN: KYROS_AUTH_SECRET is set but no longer used. The registry service now validates "+
				"bearer tokens against the Keycloak JWKS endpoint. Remove the env var.")
	}

	srv.ListenAndServe()
}
