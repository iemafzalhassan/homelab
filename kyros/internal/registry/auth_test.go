package registry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap/zaptest"
)

// stubValidator lets tests inject any token-validation behaviour without
// standing up a Keycloak. Production code uses *auth.Validator.
type stubValidator struct {
	claims map[string]any
	err    error
}

func (s *stubValidator) ValidateToken(_ string) (map[string]any, error) {
	return s.claims, s.err
}

func TestAuthMiddleware_NoToken(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	auth := NewAuthMiddleware(logger, "test-realm", &stubValidator{claims: map[string]any{"sub": "x"}})

	req := httptest.NewRequest("GET", "/v2/", nil)
	rr := httptest.NewRecorder()

	handler := auth.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	expectedAuthHeader := `Bearer realm="test-realm",service="kyros-registry"`
	if rr.Header().Get("Www-Authenticate") != expectedAuthHeader {
		t.Errorf("handler returned wrong Www-Authenticate header: got %v want %v",
			rr.Header().Get("Www-Authenticate"), expectedAuthHeader)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	auth := NewAuthMiddleware(logger, "test-realm", &stubValidator{claims: map[string]any{"sub": "user123"}})

	req := httptest.NewRequest("GET", "/v2/", nil)
	req.Header.Set("Authorization", "Bearer good-token")
	rr := httptest.NewRecorder()

	handler := auth.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the middleware stashed claims on the context.
		claims := ClaimsFromContext(r.Context())
		if claims["sub"] != "user123" {
			t.Errorf("expected claims.sub == user123, got %v", claims["sub"])
		}
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	auth := NewAuthMiddleware(logger, "test-realm", &stubValidator{err: context.DeadlineExceeded})

	req := httptest.NewRequest("GET", "/v2/", nil)
	req.Header.Set("Authorization", "Bearer bad-token")
	rr := httptest.NewRecorder()

	handler := auth.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("downstream handler must not run when token is invalid")
	}))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

// TestAuthMiddleware_NilValidator is the regression test for the
// dev-mode-bypass removal. The previous implementation accepted any token
// when the secret was empty; the new implementation must reject all requests
// when the validator is nil.
func TestAuthMiddleware_NilValidator(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	auth := NewAuthMiddleware(logger, "test-realm", nil)

	req := httptest.NewRequest("GET", "/v2/", nil)
	req.Header.Set("Authorization", "Bearer any-token")
	rr := httptest.NewRecorder()

	handler := auth.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("downstream handler must not run when validator is nil")
	}))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusServiceUnavailable {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusServiceUnavailable)
	}
}
