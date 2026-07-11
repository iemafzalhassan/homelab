package httpserver

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/kyros-platform/kyros/internal/health"
)

type Server struct {
	Router *chi.Mux
	port   string
	log    *zap.Logger
}

func New(port string, log *zap.Logger, serviceName string) *Server {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Prometheus metrics
	r.Handle("/metrics", promhttp.Handler())

	// Health endpoints
	h := &health.Handler{
		ServiceName: serviceName,
		Version:     "dev",
		Checkers:    make(health.ReadinessCheckers),
		Log:         log,
	}
	h.RegisterRoutes(r)

	return &Server{
		Router: r,
		port:   port,
		log:    log,
	}
}

func (s *Server) ListenAndServe() {
	srv := &http.Server{
		Addr:    ":" + s.port,
		Handler: s.Router,
	}

	go func() {
		s.log.Info("Starting server", zap.String("port", s.port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Fatal("Server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		s.log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	s.log.Info("Server exiting")
}
