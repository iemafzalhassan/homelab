package middleware

import (
	"context"
	"net/http"

	"github.com/kyros-platform/kyros/services/auth/internal/service"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RBACMiddleware enforces role-based access control.
type RBACMiddleware struct {
	rbacService *service.RBACService
	logger      *zap.SugaredLogger
}

// NewRBACMiddleware creates a new RBAC middleware.
func NewRBACMiddleware(rbacService *service.RBACService, logger *zap.SugaredLogger) *RBACMiddleware {
	return &RBACMiddleware{
		rbacService: rbacService,
		logger:      logger,
	}
}

// RequireRole ensures the user has the specified role.
func (m *RBACMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value("user_id").(uuid.UUID)
			if !ok {
				m.logger.Infow("user ID not found in context")
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			hasRole, err := m.rbacService.UserHasRole(r.Context(), userID, role)
			if err != nil {
				m.logger.Errorw("failed to check user role", "error", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			if !hasRole {
				m.logger.Infow("user lacks required role", "user_id", userID, "required_role", role)
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequirePermission ensures the user has the specified permission.
func (m *RBACMiddleware) RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value("user_id").(uuid.UUID)
			if !ok {
				m.logger.Infow("user ID not found in context")
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			hasPerm, err := m.rbacService.UserHasPermission(r.Context(), userID, permission)
			if err != nil {
				m.logger.Errorw("failed to check user permission", "error", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			if !hasPerm {
				m.logger.Infow("user lacks required permission", "user_id", userID, "required_permission", permission)
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}