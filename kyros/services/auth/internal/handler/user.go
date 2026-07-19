package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/kyros-platform/kyros/services/auth/internal/repository"
	"github.com/kyros-platform/kyros/services/auth/internal/service"
	"go.uber.org/zap"
)

// UserHandler handles user-related HTTP requests.
type UserHandler struct {
	userService *service.UserService
	logger      *zap.SugaredLogger
}

// NewUserHandler creates a new user handler.
func NewUserHandler(userService *service.UserService, logger *zap.SugaredLogger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// GetCurrentUser returns the currently authenticated user.
func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		h.logger.Errorw("user ID not found in context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		h.logger.Errorw("failed to get user", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, user)
}

// UpdateCurrentUser updates the currently authenticated user's profile.
func (h *UserHandler) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		h.logger.Errorw("user ID not found in context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		DisplayName string `json:"display_name"`
		AvatarURL   string `json:"avatar_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorw("invalid request body", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.UpdateProfile(r.Context(), userID, req.DisplayName, req.AvatarURL)
	if err != nil {
		h.logger.Errorw("failed to update profile", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Infow("user profile updated", "user_id", userID)
	respondJSON(w, http.StatusOK, user)
}

// ListUsers returns a paginated list of users (admin only).
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// TODO: add pagination query params
	limit := 100
	offset := 0

	users, err := h.userService.List(r.Context(), limit, offset)
	if err != nil {
		h.logger.Errorw("failed to list users", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, users)
}

// GetUserByID returns a user by ID (admin or self).
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		h.logger.Errorw("user ID not found in context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	targetID := chi.URLParam(r, "id")
	if targetID == "" {
		http.Error(w, "user ID is required", http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(targetID)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	// Check if the user is requesting their own data or is an admin (we'll implement admin check later)
	if uuid != userID {
		// For now, only allow users to view their own data
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	user, err := h.userService.GetByID(r.Context(), uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		h.logger.Errorw("failed to get user", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, user)
}

// UpdateUser updates a user by ID (admin only).
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// TODO: implement admin check
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// DeleteUser deletes a user by ID (admin only).
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// TODO: implement admin check
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// Helper function to write JSON response.
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}