package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/kyros-platform/kyros/services/auth/internal/repository"
	"github.com/kyros-platform/kyros/services/auth/internal/service"
	"go.uber.org/zap"
)

// RBACHandler handles RBAC-related HTTP requests.
type RBACHandler struct {
	rbacService *service.RBACService
	logger      *zap.SugaredLogger
}

// NewRBACHandler creates a new RBAC handler.
func NewRBACHandler(rbacService *service.RBACService, logger *zap.SugaredLogger) *RBACHandler {
	return &RBACHandler{
		rbacService: rbacService,
		logger:      logger,
	}
}

// CreateRole creates a new role.
func (h *RBACHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	// TODO: implement admin check
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
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

	role := &service.RBACService{ ... }
	 }
	}
	
	role := &repository.Role{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.rbacService.roleRepo.Create(r.Context(), role); err != nil {
		h.logger.Errorw("failed to create role", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Infow("role created", "role_id", role.ID)
	respondJSON(w, http.StatusCreated, role)
}

// ListRoles returns a list of roles.
func (h *RBACHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	// TODO: implement pagination
	limit := 100
	offset := 0

	roles, err := h.rbacService.roleRepo.List(r.Context(), limit, offset)
	if err != nil {
		h.logger.Errorw("failed to list roles", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, roles)
}

// GetRoleByID returns a role by ID.
func (h *RBACHandler) GetRoleByID(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "id")
	if roleID == "" {
		http.Error(w, "role ID is required", http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(roleID)
	if err != nil {
		http.Error(w, "invalid role ID", http.StatusBadRequest)
		return
	}

	role, err := h.rbacService.roleRepo.GetByID(r.Context(), uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "role not found", http.StatusNotFound)
			return
		}
		h.logger.Errorw("failed to get role", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, role)
}

// UpdateRole updates a role.
func (h *RBACHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	// TODO: implement admin check
	roleID := chi.URLParam(r, "id")
	if roleID == "" {
		http.Error(w, "role ID is required", http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(roleID)
	if err != nil {
		http.Error(w, "invalid role ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorw("invalid request body", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	role, err := h.rbacService.roleRepo.GetByID(r.Context(), uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "role not found", http.StatusNotFound)
			return
		}
		h.logger.Errorw("failed to get role", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	role.Name = req.Name
	role.Description = sql.NullString{String: req.Description, Valid: req.Description != ""}
	role.UpdatedAt = time.Now()

	if err := h.rbacService.roleRepo.Update(r.Context(), role); err != nil {
		h.logger.Errorw("failed to update role", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Infow("role updated", "role_id", role.ID)
	respondJSON(w, http.StatusOK, role)
}

// DeleteRole deletes a role.
func (h *RBACHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	// TODO: implement admin check
	roleID := chi.URLParam(r, "id")
	if roleID == "" {
		http.Error(w, "role ID is required", http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(roleID)
	if err != nil {
		http.Error(w, "invalid role ID", http.StatusBadRequest)
		return
	}

	if err := h.rbacService.roleRepo.Delete(r.Context(), uuid); err != nil {
		h.logger.Errorw("failed to delete role", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Infow("role deleted", "role_id", uuid)
	w.WriteHeader(http.StatusNoContent)
}

// GrantPermission grants a permission to a role.
func (h *RBACHandler) GrantPermission(w http.ResponseWriter, r *http.Request) {
	// TODO: implement admin check
	roleID := chi.URLParam(r, "roleId")
	permissionID := chi.URLParam(r, "permissionId")
	if roleID == "" || permissionID == "" {
		http.Error(w, "role ID and permission ID are required", http.StatusBadRequest)
		return
	}

	roleUUID, err := uuid.Parse(roleID)
	if err != nil {
		http.Error(w, "invalid role ID", http.StatusBadRequest)
		return
	}

	permUUID, err := uuid.Parse(permissionID)
	if err != nil {
		http.Error(w, "invalid permission ID", http.StatusBadRequest)
		return
	}

	if err := h.rbacService.GrantPermissionToRole(r.Context(), roleUUID, permUUID); err != nil {
		h.logger.Errorw("failed to grant permission", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Infow("permission granted", "role_id", roleUUID, "permission_id", permUUID)
	w.WriteHeader(http.StatusNoContent)
}

// RevokePermission revokes a permission from a role.
func (h *RBACHandler) RevokePermission(w http.ResponseWriter, r *http.Request) {
	// TODO: implement admin check
	roleID := chi.URLParam(r, "roleId")
	permissionID := chi.URLParam(r, "permissionId")
	if roleID == "" || permissionID == "" {
		http.Error(w, "role ID and permission ID are required", http.StatusBadRequest)
		return
	}

	roleUUID, err := uuid.Parse(roleID)
	if err != nil {
		http.Error(w, "invalid role ID", http.StatusBadRequest)
		return
	}

	permUUID, err := uuid.Parse(permissionID)
	if err != nil {
		http.Error(w, "invalid permission ID", http.StatusBadRequest)
		return
	}

	if err := h.rbacService.RevokePermissionFromRole(r.Context(), roleUUID, permUUID); err != nil {
		h.logger.Errorw("failed to revoke permission", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Infow("permission revoked", "role_id", roleUUID, "permission_id", permUUID)
	w.WriteHeader(http.StatusNoContent)
}

// AssignRole assigns a role to a user.
func (h *RBACHandler) AssignRole(w http.ResponseWriter, r *http.Request) {
	// TODO: implement admin check
	userID := chi.URLParam(r, "userId")
	roleID := chi.URLParam(r, "roleId")
	if userID == "" || roleID == "" {
		http.Error(w, "user ID and role ID are required", http.StatusBadRequest)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	roleUUID, err := uuid.Parse(roleID)
	if err != nil {
		http.Error(w, "invalid role ID", http.StatusBadRequest)
		return
	}

	if err := h.rbacService.AssignRoleToUser(r.Context(), userUUID, roleUUID); err != nil {
		h.logger.Errorw("failed to assign role", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Infow("role assigned", "user_id", userUUID, "role_id", roleUUID)
	w.WriteHeader(http.StatusNoContent)
}

// RevokeRole removes a role from a user.
func (h *RBACHandler) RevokeRole(w http.ResponseWriter, r *http.Request) {
	// TODO: implement admin check
	userID := chi.URLParam(r, "userId")
	roleID := chi.URLParam(r, "roleId")
	if userID == "" || roleID == "" {
		http.Error(w, "user ID and role ID are required", http.StatusBadRequest)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	roleUUID, err := uuid.Parse(roleID)
	if err != nil {
		http.Error(w, "invalid role ID", http.StatusBadRequest)
		return
	}

	if err := h.rbacService.RevokeRoleFromUser(r.Context(), userUUID, roleUUID); err != nil {
		h.logger.Errorw("failed to revoke role", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Infow("role revoked", "user_id", userUUID, "role_id", roleUUID)
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