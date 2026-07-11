package health

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Checker interface {
	Check(ctx context.Context) error
}

type ReadinessCheckers map[string]Checker

type Handler struct {
	ServiceName string
	Version     string
	Checkers    ReadinessCheckers
	Log         *zap.Logger
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/health", h.Health)
	r.Get("/live", h.Live)
	r.Get("/ready", h.Ready)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"` + h.ServiceName + `","version":"` + h.Version + `"}`))
}

func (h *Handler) Live(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	for name, checker := range h.Checkers {
		if err := checker.Check(r.Context()); err != nil {
			h.Log.Error("Readiness check failed", zap.String("checker", name), zap.Error(err))
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("unavailable"))
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}
