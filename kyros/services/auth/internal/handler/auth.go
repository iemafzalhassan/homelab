package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/kyros-platform/kyros/services/auth/internal/service"
	"github.com/kyros-platform/kyros/services/auth/internal/repository"
	"go.uber.org/zap"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	authService *service.AuthService
	userService *service.UserService
	logger      *zap.SugaredLogger
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(authService *service.AuthService, userService *service.UserService, logger *zap.SugaredLogger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
		logger:      logger,
	}
}

// Register handles user registration.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email        string `json:"email"`
		DisplayName  string `json:"display_name"`
		Password     string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorw("invalid request body", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" || req.DisplayName == "" {
		http.Error(w, "email, display_name, and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Register(r.Context(), req.Email, req.DisplayName, req.Password)
	if err != nil {
		if err.Error() == "user already exists" {
			h.logger.Infow("registration failed: user already exists", "email", req.Email)
			http.Error(w, "user already exists", http.StatusConflict)
			return
		}
		h.logger.Errorw("registration failed", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Infow("user registered", "email", req.Email, "user_id", user.ID)
	respondJSON(w, http.StatusCreated, user)
}

// Login handles user login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorw("invalid request body", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, user, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			h.logger.Infow("login failed: invalid credentials", "email", req.Email)
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		h.logger.Errorw("login failed", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Set refresh token as HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // should be true in production
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(30 * 24 * time.Hour), // 30 days
	})

	h.logger.Infow("user logged in", "email", req.Email, "user_id", user.ID)
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"access_token":  accessToken,
		"token_type":    "Bearer",
		"expires_in":    int64(h.authService.accessTokenExp.Seconds()),
		"user":          user,
	})
}

// Logout handles user logout by invalidating the refresh token.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			h.logger.Infow("logout attempt without refresh token")
			http.Error(w, "no refresh token", http.StatusBadRequest)
			return
		}
		h.logger.Errorw("failed to get cookie", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	refreshToken := cookie.Value
	session, err := h.authService.sessionRepo.GetByRefreshToken(r.Context(), refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.logger.Infow("logout: session not found", "refresh_token", refreshToken)
			// Still clear the cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "refresh_token",
				Value:    "",
				Path:     "/",
				HttpOnly: true,
				Expires:  time.Now().Add(-time.Hour),
			})
			w.WriteHeader(http.StatusOK)
			return
		}
		h.logger.Errorw("logout: failed to get session", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := h.authService.sessionRepo.Delete(r.Context(), session.ID); err != nil {
		h.logger.Errorw("logout: failed to delete session", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Clear the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(-time.Hour),
	})

	h.logger.Infow("user logged out", "user_id", session.UserID)
	w.WriteHeader(http.StatusOK)
}

// RefreshToken handles refreshing the access token using the refresh token cookie.
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			h.logger.Infow("refresh token missing")
			http.Error(w, "refresh token required", http.StatusUnauthorized)
			return
		}
		h.logger.Errorw("failed to get cookie", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	refreshToken := cookie.Value
	session, err := h.authService.sessionRepo.GetByRefreshToken(r.Context(), refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.logger.Infow("refresh token not found or expired", "token", refreshToken)
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		}
		h.logger.Errorw("failed to get session", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		h.logger.Infow("refresh token expired", "user_id", session.UserID)
		http.Error(w, "refresh token expired", http.StatusUnauthorized)
		return
	}

	// Generate new access token
	user, err := h.authService.userRepo.GetByID(r.Context(), session.UserID)
	if err != nil {
		h.logger.Errorw("failed to get user for refresh", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	accessToken, err := h.authService.generateAccessToken(user)
	if err != nil {
		h.logger.Errorw("failed to generate access token", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Infow("token refreshed", "user_id", user.ID)
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"access_token":  accessToken,
		"token_type":    "Bearer",
		"expires_in":    int64(h.authService.accessTokenExp.Seconds()),
	})
}

// CreateAPIKey handles creating a new API key for the authenticated user.
func (h *AuthHandler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	// The user ID is obtained from the auth middleware (we'll implement that later)
	// For now, we'll get it from context (key "user_id")
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		h.logger.Errorw("user ID not found in context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorw("invalid request body", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	// Generate a random API key
	apiKey := make([]byte, 32)
	if _, err := rand.Read(apiKey); err != nil {
		h.logger.Errorw("failed to generate random bytes", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	apiKeyString := base64.RawURLEncoding.EncodeToString(apiKey)

	// Hash the API key
	hashedKey, err := bcrypt.GenerateFromPassword([]byte(apiKeyString), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Errorw("failed to hash API key", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Create the API key record
	key := &repository.APIKey{
		ID:        uuid.New(),
		UserID:    userID,
		KeyHash:   string(hashedKey),
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.authService.apiKeyRepo.Create(r.Context(), key); err != nil {
		h.logger.Errorw("failed to create API key", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Infow("API key created", "user_id", userID, "key_id", key.ID)
	// Return the plaintext key only once
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"id":        key.ID,
		"name":      key.Name,
		"key":       apiKeyString, // only shown once
		"created_at": key.CreatedAt,
	})
}

// DeleteAPIKey deletes an API key by ID.
func (h *AuthHandler) DeleteAPIKey(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		h.logger.Errorw("user ID not found in context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	keyID := chi.URLParam(r, "id")
	if keyID == "" {
		http.Error(w, "key ID is required", http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(keyID)
	if err != nil {
		http.Error(w, "invalid key ID", http.StatusBadRequest)
		return
	}

	// Verify the key belongs to the user
	key, err := h.authService.apiKeyRepo.GetByID(r.Context(), uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "API key not found", http.StatusNotFound)
			return
		}
		h.logger.Errorw("failed to get API key", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if key.UserID != userID {
		h.logger.Infow("attempt to delete another user's API key", "user_id", userID, "key_id", uuid)
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if err := h.authService.apiKeyRepo.Delete(r.Context(), uuid); err != nil {
		h.logger.Errorw("failed to delete API key", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Infow("API key deleted", "user_id", userID, "key_id", uuid)
	w.WriteHeader(http.StatusNoContent)
}

// Helper function to write JSON response.
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}