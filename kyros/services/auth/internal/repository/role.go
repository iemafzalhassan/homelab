package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Role represents a role.
type Role struct {
	ID          uuid.UUID
	Name        string
	Description sql.NullString
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// RoleRepository handles storage of roles.
type RoleRepository struct {
	db *sql.DB
}

// NewRoleRepository creates a new role repository.
func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

// Create inserts a new role.
func (r *RoleRepository) Create(ctx context.Context, role *Role) error {
	query := `
		INSERT INTO roles (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`
	row := r.db.QueryRowContext(ctx, query,
		role.ID,
		role.Name,
		role.Description,
		role.CreatedAt,
		role.UpdatedAt,
	)
	return row.Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt)
}

// GetByName retrieves a role by name.
func (r *RoleRepository) GetByName(ctx context.Context, name string) (*Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		WHERE name = $1`
	row := r.db.QueryRowContext(ctx, query, name)

	var role Role
	err := row.Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetByID retrieves a role by ID.
func (r *RoleRepository) GetByID(ctx context.Context, id uuid.UUID) (*Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var role Role
	err := row.Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// List returns a list of roles with pagination.
func (r *RoleRepository) List(ctx context.Context, limit, offset int) ([]*Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		ORDER BY name
		LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*Role
	for rows.Next() {
		var role Role
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

// Update updates a role.
func (r *RoleRepository) Update(ctx context.Context, role *Role) error {
	query := `
		UPDATE roles
		SET name = $2, description = $3, updated_at = NOW()
		WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query,
		role.ID,
		role.Name,
		role.Description,
	)
	return err
}

// Delete removes a role.
func (r *RoleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM roles WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}