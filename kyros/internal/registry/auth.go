package registry

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// claimsKey is the context key under which validated token claims are stored.
// Handlers downstream of RequireAuth can pull claims with
// `registry.ClaimsFromContext(r.Context())`.
type claimsKey struct{}

// ClaimsFromContext returns the JWT claims placed on the request context by
// RequireAuth, or nil if no auth middleware ran on this request.
func ClaimsFromContext(ctx context.Context) map[string]any {
	v, _ := ctx.Value(claimsKey{}).(map[string]any)
	return v
}

// TokenValidator validates a bearer token and returns its claims on success.
// The production implementation is *auth.Validator (Keycloak JWKS, RS256).
// Tests can pass a stub.
type TokenValidator interface {
	ValidateToken(tokenString string) (claims map[string]any, err error)
}

// AuthMiddleware guards OCI registry endpoints. It enforces a valid bearer
// token via the supplied TokenValidator and exposes the validated claims on
// the request context.
//
// Security: the previous version of this middleware accepted any token when
// the configured secret was empty ("dev mode"). That was a footgun: an
// operator who forgot to set the secret in production would silently allow
// every push. The middleware now requires a non-nil validator. If a
// deployment wants to skip auth (e.g. local dev against a backend registry
// with its own auth), it should configure the backend registry to require
// auth and the validator to validate against the dev Keycloak — NOT skip
// the middleware.
type AuthMiddleware struct {
	logger    *zap.SugaredLogger
	realm     string
	validator TokenValidator
}

func NewAuthMiddleware(logger *zap.SugaredLogger, realm string, validator TokenValidator) *AuthMiddleware {
	if validator == nil {
		// Fail-closed. Returning a middleware that 503s on every request is
		// far safer than one that allows every request.
		logger.Errorw("AuthMiddleware constructed with nil validator; registry will reject all requests")
	}
	return &AuthMiddleware{
		logger:    logger,
		realm:     realm,
		validator: validator,
	}
}

func (a *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Refuse early if the middleware is misconfigured.
		if a.validator == nil {
			a.logger.Errorw("refusing request: auth middleware has no validator")
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"errors":[{"code":"UNAUTHORIZED","message":"registry not configured","detail":null}]}`))
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			a.challenge(w, r)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			a.challenge(w, r)
			return
		}

		claims, err := a.validator.ValidateToken(parts[1])
		if err != nil {
			a.logger.Warnw("token validation failed", "error", err)
			a.challenge(w, r)
			return
		}

		// Stash claims for downstream handlers (logging, rate limiting, future
		// per-user metrics).
		ctx := context.WithValue(r.Context(), claimsKey{}, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *AuthMiddleware) challenge(w http.ResponseWriter, r *http.Request) {
	service := "kyros-registry"
	scope := r.URL.Query().Get("scope")

	challenge := fmt.Sprintf(`Bearer realm="%s",service="%s"`, a.realm, service)
	if scope != "" {
		challenge += fmt.Sprintf(`,scope="%s"`, scope)
	}

	w.Header().Set("Www-Authenticate", challenge)
	w.WriteHeader(http.StatusUnauthorized)

	// Docker client expects a JSON error response.
	_, _ = w.Write([]byte(`{"errors":[{"code":"UNAUTHORIZED","message":"authentication required","detail":null}]}`))
}
