package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/kyros-platform/kyros/services/auth/internal/repository"
	"github.com/kyros-platform/kyros/services/auth/internal/service"
	"go.uber.org/zap"
)

// OrganizationHandler handles organization-related HTTP requests.
type OrganizationHandler struct {
	orgService *service.OrganizationService
	logger     *zap.SugaredLogger
}

// NewOrganizationHandler creates a new organization handler.
func NewOrganizationHandler(orgService *service.OrganizationService, logger *zap.SugaredLogger) *OrganizationHandler {
	return &OrganizationHandler{
		orgService: orgService,
		logger:     logger,
	}
}

// CreateOrganization creates a new organization.
func (h *OrganizationHandler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (authenticated user)
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		h.logger.Errorw("user ID not found in context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Slug        string `json:"slug"`
		DisplayName string `json:"display_name"`
		Description string `json:"description"`
		AvatarURL   string `json:"avatar_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorw("invalid request body", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Slug == "" || req.DisplayName == "" {
		http.Error(w, "slug and display_name are required", http.StatusBadRequest)
		return
	}

	org, err := h.orgService.Create(r.Context(), req.Slug, req.DisplayName, req.Description, req.AvatarURL)
	if err != nil {
		h.logger.Errorw("failed to create organization", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Infow("organization created", "user_id", userID, "org_id", org.ID)
	respondJSON(w, http.StatusCreated, org)
}

// ListOrganizations returns a paginated list of organizations.
func (h *OrganizationHandler) ListOrganizations(w http.ResponseWriter, r *http.Request) {
	// TODO: add pagination
	limit := 100
	offset := 0

	orgs, err := h.orgService.List(r.Context(), limit, offset)
	if err != nil {
		h.logger.Errorw("failed to list organizations", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, orgs)
}

// GetOrganizationByID returns an organization by ID.
func (h *OrganizationHandler) GetOrganizationByID(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "id")
	if orgID == "" {
		http.Error(w, "organization ID is required", http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(orgID)
	if err != nil {
		http.Error(w, "invalid organization ID", http.StatusBadRequest)
		return
	}

	org, err := h.orgService.GetByID(r.Context(), uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "organization not found", http.StatusNotFound)
			return
		}
		h.logger.Errorw("failed to get organization", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, org)
}

// UpdateOrganization updates an organization.
func (h *OrganizationHandler) UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "id")
	if orgID == "" {
		http.Error(w, "organization ID is required", http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(orgID)
	if err != nil {
		http.Error(w, "invalid organization ID", http.StatusBadRequest)
		return
	}

	// Get user ID from context (authenticated user)
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		h.logger.Errorw("user ID not found in context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Slug        string `json:"slug"`
		DisplayName string `json:"display_name"`
		Description string `json:"description"`
		AvatarURL   string `json:"avatar_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorw("invalid request body", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	org, err := h.orgService.Update(r.Context(), uuid, req.Slug, req.DisplayName, req.Description, req.AvatarURL)
	if err != nil {
		h.logger.Errorw("failed to update organization", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Infow("organization updated", "user_id", userID, "org_id", org.ID)
	respondJSON(w, http.StatusOK, org)
}

// DeleteOrganization deletes an organization.
func (h *OrganizationHandler) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "id")
	if orgID == "" {
		http.Error(w, "organization ID is required", http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(orgID)
	if err != nil {
		http.Error(w, "invalid organization ID", http.StatusBadRequest)
		return
	}

	// Get user ID from context (authenticated user)
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		h.logger.Errorw("user ID not found in context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.orgService.Delete(r.Context(), uuid); err != nil {
		h.logger.Errorw("failed to delete organization", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Infow("organization deleted", "user_id", userID, "org_id", uuid)
	w.WriteHeader(http.StatusNoContent)
}

// AddMember adds a user to an organization.
func (h *OrganizationHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// RemoveMember removes a user from an organization.
func (h *OrganizationHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
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