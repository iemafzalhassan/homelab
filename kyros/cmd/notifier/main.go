package main

import (
	"github.com/kyros-platform/kyros/internal/config"
	"github.com/kyros-platform/kyros/internal/httpserver"
	"github.com/kyros-platform/kyros/internal/logger"
)

func main() {
	cfg := config.Load("notifier")
	log := logger.New(cfg.LogLevel, cfg.Env)
	defer log.Sync()

	srv := httpserver.New(cfg.Port, log, "notifier")

	// Register service-specific routes here

	srv.ListenAndServe()
}
