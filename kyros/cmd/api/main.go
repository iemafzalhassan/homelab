package main

import (
	"os"

	"github.com/kyros-platform/kyros/internal/api/auth"
	"github.com/kyros-platform/kyros/internal/api/webhooks"
	"github.com/kyros-platform/kyros/internal/config"
	"github.com/kyros-platform/kyros/internal/httpserver"
	"github.com/kyros-platform/kyros/internal/logger"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load("api")
	log := logger.New(cfg.LogLevel, cfg.Env)
	defer func() { _ = log.Sync() }()

	srv := httpserver.New(cfg.Port, log, "api")

	webhookHandler, err := webhooks.NewRegistryHandler(log.Sugar(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to initialize webhook handler", zap.Error(err))
	}
	defer func() { webhookHandler.Close() }()

	issuerURL := os.Getenv("KEYCLOAK_ISSUER")
	if issuerURL == "" {
		issuerURL = "http://localhost:8081/realms/kyros"
	}

	authHandler, err := auth.NewHandler(log.Sugar(), cfg.DatabaseURL, issuerURL)
	if err != nil {
		log.Fatal("failed to initialize auth handler", zap.Error(err))
	}
	defer func() { _ = authHandler.Close() }()

	srv.Router.Post("/webhooks/registry", webhookHandler.HandleEvent)
	srv.Router.Post("/v1/auth/sync", authHandler.HandleSync)

	srv.ListenAndServe()
}
