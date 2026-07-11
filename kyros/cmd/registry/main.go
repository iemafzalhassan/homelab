package main

import (
	"github.com/kyros-platform/kyros/internal/config"
	"github.com/kyros-platform/kyros/internal/httpserver"
	"github.com/kyros-platform/kyros/internal/logger"
)

func main() {
	cfg := config.Load("registry")
	log := logger.New(cfg.LogLevel, cfg.Env)
	defer log.Sync()

	srv := httpserver.New(cfg.Port, log, "registry")

	// Register service-specific routes here

	srv.ListenAndServe()
}
