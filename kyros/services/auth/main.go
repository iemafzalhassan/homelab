package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kyros-platform/kyros/internal/logger"
	"github.com/kyros-platform/kyros/internal/metrics"
	"github.com/kyros-platform/kyros/services/auth/internal/config"
	"github.com/kyros-platform/kyros/services/auth/internal/handler"
	"github.com/kyros-platform/kyros/services/auth/internal/middleware"
	"github.com/kyros-platform/kyros/services/auth/internal/repository"
	"github.com/kyros-platform/kyros/services/auth/internal/service"
	"github.com/kyros-platform/kyros/services/auth/internal/telemetry"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Load configuration
	cfg := config.Load()
	if cfg == nil {
		panic("failed to load configuration")
	}

	// Initialize logger
	log := logger.New(cfg.LogLevel, cfg.Env)
	sugar := log.Sugar()
	defer func() {
		_ = log.Sync()
	}()

	// Initialize telemetry (OpenTelemetry)
	tp, err := telemetry.InitTelemetry(cfg.ServiceName, cfg.OTELCollectorURL)
	if err != nil {
		sugar.Fatalw("failed to initialize telemetry", "error", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			sugar.Errorw("error shutting down tracer provider", "error", err)
		}
	}()

	// Prometheus metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	// Initialize database connection
	db, err := repository.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		sugar.Fatalw("failed to connect to database", "error", err)
	}
	defer db.Close()

	// Run migrations
	if err := repository.Migrate(db.DB); err != nil {
		sugar.Fatalw("failed to run migrations", "error", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	orgRepo := repository.NewOrganizationRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	apiKeyRepo := repository.NewAPIKeyRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo, sugar)
	orgService := service.NewOrganizationService(orgRepo, sugar)
	authService := service.NewAuthService(userRepo, sessionRepo, apiKeyRepo, sugar, cfg.JWTSecret, cfg.AccessTokenExpiry, cfg.RefreshTokenExpiry, cfg.BCryptCost)
	rbacService := service.NewRBACService(roleRepo, permissionRepo, userRepo, sugar)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, sugar)
	userHandler := handler.NewUserHandler(userService, sugar)
	orgHandler := handler.NewOrganizationHandler(orgService, sugar)
	rbacHandler := handler.NewRBACHandler(rbacService, sugar)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(sugar, cfg.JWTSecret)
	rbacMiddleware := middleware.NewRBACMiddleware(rbacService, sugar)
	// permissionMiddleware := middleware.NewPermissionMiddleware(rbacService, sugar) // TODO: use if needed

	// Set up HTTP router
	r := chi.NewRouter()
	// Middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	// Custom middleware
	r.Use(func(next http.Handler) http.Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			// Add request ID to context if not present
			ctx := r.Context()
			if ctx.Value(middleware.RequestIDKey) == nil {
				ctx = context.WithValue(ctx, middleware.RequestIDKey, middleware.GetReqID(ctx))
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
	r.Use(func(next http.Handler) http.Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			// Log request
			sugar.Infow("request started",
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
				"request_id", r.Context().Value(middleware.RequestIDKey),
			)
			start := time.Now()
			rw := &middleware.ResponseWriter{ResponseWriter: w, Status: 200}
			next.ServeHTTP(rw, r)
			// Log response
			sugar.Infow("request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.Status,
				"duration_ms", time.Since(start).Milliseconds(),
				"request_id", r.Context().Value(middleware.RequestIDKey),
			)
			// Update Prometheus metrics
			metrics.HttpRequestsTotal.WithLabelValues("auth", r.Method, r.URL.Path, string(rw.Status)).Inc()
			metrics.HttpRequestDuration.WithLabelValues("auth", r.Method, r.URL.Path, string(rw.Status)).Observe(time.Since(start).Seconds())
		}
	})

	// Public routes
	r.Post("/api/v1/auth/register", authHandler.Register)
	r.Post("/api/v1/auth/login", authHandler.Login)
	r.Post("/api/v1/auth/refresh", authHandler.RefreshToken)
	r.Post("/api/v1/auth/logout", authHandler.Logout)
	r.Post("/api/v1/auth/api-keys", authHandler.CreateAPIKey, authMiddleware)
	r.Delete("/api/v1/auth/api-keys/{id}", authHandler.DeleteAPIKey, authMiddleware)

	// Protected routes (require authentication)
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Get("/api/v1/users/me", userHandler.GetCurrentUser)
		r.Put("/api/v1/users/me", userHandler.UpdateCurrentUser)
		r.Get("/api/v1/users", userHandler.ListUsers)
		r.Get("/api/v1/users/{id}", userHandler.GetUserByID)
		r.Put("/api/v1/users/{id}", userHandler.UpdateUser)
		r.Delete("/api/v1/users/{id}", userHandler.DeleteUser)
	})

	// Organization routes (require authentication)
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Post("/api/v1/organizations", orgHandler.CreateOrganization)
		r.Get("/api/v1/organizations", orgHandler.ListOrganizations)
		r.Get("/api/v1/organizations/{id}", orgHandler.GetOrganizationByID)
		r.Put("/api/v1/organizations/{id}", orgHandler.UpdateOrganization)
		r.Delete("/api/v1/organizations/{id}", orgHandler.DeleteOrganization)
		r.Post("/api/v1/organizations/{id}/members", orgHandler.AddMember)
		r.Delete("/api/v1/organizations/{id}/members/{userId}", orgHandler.RemoveMember)
	})

	// RBAC routes (require authentication and admin role)
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Use(rbacMiddleware.RequireRole("admin"))
		r.Post("/api/v1/rbac/roles", rbacHandler.CreateRole)
		r.Get("/api/v1/rbac/roles", rbacHandler.ListRoles)
		r.Get("/api/v1/rbac/roles/{id}", rbacHandler.GetRoleByID)
		r.Put("/api/v1/rbac/roles/{id}", rbacHandler.UpdateRole)
		r.Delete("/api/v1/rbac/roles/{id}", rbacHandler.DeleteRole)
		r.Post("/api/v1/rbac/roles/{roleId}/permissions/{permissionId}", rbacHandler.GrantPermission)
		r.Delete("/api/v1/rbac/roles/{roleId}/permissions/{permissionId}", rbacHandler.RevokePermission)
		r.Post("/api/v1/rbac/users/{userId}/roles/{roleId}", rbacHandler.AssignRoleToUser)
		r.Delete("/api/v1/rbac/users/{userId}/roles/{roleId}", rbacHandler.RevokeRoleFromUser)
	})

	// Health endpoints
	r.Get("/healthz", handler.HealthHandler)
	r.Get("/readyz", handler.ReadinessHandler(db))
	r.Get("/livez", handler.LivenessHandler)

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		sugar.Infow("starting server", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Fatalw("server failed to start", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	sugar.Infow("shutting down server")

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		sugar.Fatalw("server forced to shutdown", "error", err)
	}

	sugar.Infow("server exited")
}