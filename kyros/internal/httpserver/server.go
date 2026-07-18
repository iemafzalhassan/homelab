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
	// NOTE: middleware.RealIP was removed. It is deprecated in go-chi v5
	// because it unconditionally trusts X-Forwarded-For / X-Real-IP /
	// True-Client-IP — which is exploitable when the service is reachable
	// without going through a known proxy (e.g. a debug port-forward
	// during a CVE triage). RemoteAddr remains the source of truth; any
	// proxy-header trust must be wired explicitly at the edge (Traefik
	// middleware) and passed in via a configured IP extractor.
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
	// ReadHeaderTimeout caps the time a client may spend sending the
	// request line + headers. Without it, a slow/large-header client can
	// hold a goroutine indefinitely (Slowloris — gosec G112). We set it
	// to 5s: long enough for any realistic client, short enough to bound
	// the resource a single misbehaving connection can consume.
	//
	// ReadTimeout is intentionally NOT set: the per-request Timeout
	// middleware already bounds the body read + handler execution, and
	// services that accept large request bodies (the OCI registry blob
	// push endpoint) need ReadTimeout=0.
	srv := &http.Server{
		Addr:              ":" + s.port,
		Handler:           s.Router,
		ReadHeaderTimeout: 5 * time.Second,
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
