package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/kyros-platform/kyros/services/auth/internal/service"
	"go.uber.org/zap"
)

// AuthMiddleware validates the access token and sets the user ID in the context.
type AuthMiddleware struct {
	authService *service.AuthService
	logger      *zap.SugaredLogger
}

// NewAuthMiddleware creates a new auth middleware.
func NewAuthMiddleware(authService *service.AuthService, logger *zap.SugaredLogger) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		logger:      logger,
	}
}

// ValidateToken validates the access token from the Authorization header.
func (m *AuthMiddleware) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.logger.Infow("missing authorization header")
			http.Error(w, "authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			m.logger.Infow("invalid authorization header format")
			http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			m.logger.Infow("invalid token", "error", err)
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// Extract user ID from token claims (sub)
		userIDStr, ok := claims["sub"].(string)
		if !ok {
			m.logger.Infow("token missing sub claim")
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			m.logger.Infow("invalid user ID in token", "error", err)
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// Add user ID to context
	ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}