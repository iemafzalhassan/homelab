package service

import (
	"context"
	"errors"

	"github.com/kyros-platform/kyros/services/auth/internal/repository"
)

// RBACService handles role-based access control.
type RBACService struct {
	roleRepo      *repository.RoleRepository
	permissionRepo *repository.PermissionRepository
	userRepo      *repository.UserRepository
	logger        *repository.Logger
}

// NewRBACService creates a new RBAC service.
func NewRBACService(roleRepo *repository.RoleRepository, permissionRepo *repository.PermissionRepository, userRepo *repository.UserRepository, logger *repository.Logger) *RBACService {
	return &RBACService{
		roleRepo:      roleRepo,
		permissionRepo: permissionRepo,
		userRepo:      userRepo,
		logger:        logger,
	}
}

// AssignRoleToUser assigns a role to a user.
func (s *RBACService) AssignRoleToUser(ctx context.Context, userID, roleID uuid.UUID) error {
	// Check if user exists
	if _, err := s.userRepo.GetByID(ctx, userID); err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	// Check if role exists
	if _, err := s.roleRepo.GetByID(ctx, roleID); err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	// Insert into user_roles table
	query := `
		INSERT INTO user_roles (user_id, role_id, assigned_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (user_id, role_id) DO NOTHING`
	_, err := s.userRepo.db.ExecContext(ctx, query, userID, roleID)
	return err
}

// RevokeRoleFromUser removes a role from a user.
func (s *RBACService) RevokeRoleFromUser(ctx context.Context, userID, roleID uuid.UUID) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`
	_, err := s.userRepo.db.ExecContext(ctx, query, userID, roleID)
	return err
}

// GetUserRoles returns the roles assigned to a user.
func (s *RBACService) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]*repository.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at
		FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1
		ORDER BY r.name`
	rows, err := s.userRepo.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*repository.Role
	for rows.Next() {
		var role repository.Role
		if err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.CreatedAt,
			&role.UpdatedAt,
		); err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}
	return roles, rows.Err()
}

// GrantPermissionToRole grants a permission to a role.
func (s *RBACService) GrantPermissionToRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	// Check role exists
	if _, err := s.roleRepo.GetByID(ctx, roleID); err != nil {
		return fmt.Errorf("role not found: %w", err)
	}
	// Check permission exists
	if _, err := s.permissionRepo.GetByID(ctx, permissionID); err != nil {
		return fmt.Errorf("permission not found: %w", err)
	}

	query := `
		INSERT INTO role_permissions (role_id, permission_id, granted_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (role_id, permission_id) DO NOTHING`
	_, err := s.roleRepo.db.ExecContext(ctx, query, roleID, permissionID)
	return err
}

// RevokePermissionFromRole revokes a permission from a role.
func (s *RBACService) RevokePermissionFromRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	query := `DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2`
	_, err := s.roleRepo.db.ExecContext(ctx, query, roleID, permissionID)
	return err
}

// GetRolePermissions returns the permissions for a role.
func (s *RBACService) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]*repository.Permission, error) {
	query := `
		SELECT p.id, p.name, p.description, p.created_at, p.updated_at
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
		ORDER BY p.name`
	rows, err := s.roleRepo.db.QueryContext(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*repository.Permission
	for rows.Next() {
		var perm repository.Permission
		if err := rows.Scan(
			&perm.ID,
			&perm.Name,
			&perm.Description,
			&perm.CreatedAt,
			&perm.UpdatedAt,
		); err != nil {
			return nil, err
		}
		permissions = append(permissions, &perm)
	}
	return permissions, rows.Err()
}

// UserHasPermission checks if a user has a specific permission (directly or via role).
func (s *RBACService) UserHasPermission(ctx context.Context, userID uuid.UUID, permissionName string) (bool, error) {
	// First, check if the user has the permission directly? We don't have a direct user_permissions table, so only via roles.
	// We'll check via roles.
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM user_roles ur
			JOIN role_permissions rp ON ur.role_id = rp.role_id
			JOIN permissions p ON rp.permission_id = p.id
			WHERE ur.user_id = $1 AND p.name = $2
		)`
	var exists bool
	err := s.userRepo.db.QueryRowContext(ctx, query, userID, permissionName).Scan(&exists)
	return exists, err
}

// UserHasRole checks if a user has a specific role.
func (s *RBACService) UserHasRole(ctx context.Context, userID uuid.UUID, roleName string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM user_roles ur
			JOIN roles r ON ur.role_id = r.id
			WHERE ur.user_id = $1 AND r.name = $2
		)`
	var exists bool
	err := s.userRepo.db.QueryRowContext(ctx, query, userID, roleName).Scan(&exists)
	return exists, err
}