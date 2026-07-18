package auth

import (
	"database/sql"
	"net/http"
	"strings"

	kyrosauth "github.com/kyros-platform/kyros/internal/auth"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// Handler is the HTTP handler for /v1/auth/sync. It validates the Keycloak
// access token (so we know the caller is who they say they are), then upserts
// the user into the local database. This is the source of truth for the
// user table — Keycloak owns identity, this handler owns provisioning.
type Handler struct {
	log       *zap.SugaredLogger
	db        *sql.DB
	validator *kyrosauth.Validator
}

func NewHandler(log *zap.SugaredLogger, dbURL string, issuerURL string) (*Handler, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Handler{
		log:       log,
		db:        db,
		validator: kyrosauth.NewValidator(issuerURL),
	}, nil
}

func (h *Handler) Close() error {
	return h.db.Close()
}

func (h *Handler) HandleSync(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "missing or invalid authorization header", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := h.validator.ValidateToken(tokenString)
	if err != nil {
		h.log.Warnw("failed to validate token", "error", err)
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	sub, _ := claims["sub"].(string)
	email, _ := claims["email"].(string)
	name, _ := claims["name"].(string)

	if name == "" {
		firstName, _ := claims["given_name"].(string)
		lastName, _ := claims["family_name"].(string)
		name = strings.TrimSpace(firstName + " " + lastName)
	}

	if sub == "" || email == "" {
		h.log.Warnw("token missing required claims", "sub", sub, "email", email)
		http.Error(w, "missing required claims", http.StatusBadRequest)
		return
	}

	// Upsert user into database. The "ON CONFLICT (keycloak_sub)" clause is
	// safe because the column is UNIQUE NOT NULL — see migration 000001.
	query := `
		INSERT INTO users (keycloak_sub, email, display_name)
		VALUES ($1, $2, $3)
		ON CONFLICT (keycloak_sub)
		DO UPDATE SET
			email = EXCLUDED.email,
			display_name = EXCLUDED.display_name,
			updated_at = NOW()
	`

	_, err = h.db.Exec(query, sub, email, name)
	if err != nil {
		h.log.Errorw("failed to upsert user", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.log.Infow("user synced successfully", "sub", sub, "email", email)
	w.WriteHeader(http.StatusOK)
}
