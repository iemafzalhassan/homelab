package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Permission represents a permission.
type Permission struct {
	ID          uuid.UUID
	Name        string
	Description sql.NullString
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// PermissionRepository handles storage of permissions.
type PermissionRepository struct {
	db *sql.DB
}

// NewPermissionRepository creates a new permission repository.
func NewPermissionRepository(db *sql.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

// Create inserts a new permission.
func (r *PermissionRepository) Create(ctx context.Context, p *Permission) error {
	query := `
		INSERT INTO permissions (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`
	row := r.db.QueryRowContext(ctx, query,
		p.ID,
		p.Name,
		p.Description,
		p.CreatedAt,
		p.UpdatedAt,
	)
	return row.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

// GetByName retrieves a permission by name.
func (r *PermissionRepository) GetByName(ctx context.Context, name string) (*Permission, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM permissions
		WHERE name = $1`
	row := r.db.QueryRowContext(ctx, query, name)

	var p Permission
	err := row.Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// GetByID retrieves a permission by ID.
func (r *PermissionRepository) GetByID(ctx context.Context, id uuid.UUID) (*Permission, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM permissions
		WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var p Permission
	err := row.Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// List returns a list of permissions with pagination.
func (r *PermissionRepository) List(ctx context.Context, limit, offset int) ([]*Permission, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM permissions
		ORDER BY name
		LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*Permission
	for rows.Next() {
		var p Permission
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		permissions = append(permissions, &p)
	}
	return permissions, rows.Err()
}

// Update updates a permission.
func (r *PermissionRepository) Update(ctx context.Context, p *Permission) error {
	query := `
		UPDATE permissions
		SET name = $2, description = $3, updated_at = NOW()
		WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query,
		p.ID,
		p.Name,
		p.Description,
	)
	return err
}

// Delete removes a permission.
func (r *PermissionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM permissions WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}