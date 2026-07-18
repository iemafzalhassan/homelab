package main

import (
	"github.com/kyros-platform/kyros/internal/config"
	"github.com/kyros-platform/kyros/internal/httpserver"
	"github.com/kyros-platform/kyros/internal/logger"
)

func main() {
	cfg := config.Load("sbom")
	log := logger.New(cfg.LogLevel, cfg.Env)
	defer func() { _ = log.Sync() }()

	srv := httpserver.New(cfg.Port, log, "sbom")

	// Register service-specific routes here

	srv.ListenAndServe()
}
