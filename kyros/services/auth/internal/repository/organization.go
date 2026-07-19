package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Organization represents an organization.
type Organization struct {
	ID          uuid.UUID
	Slug        string
	DisplayName string
	Description sql.NullString
	AvatarURL   sql.NullString
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// OrganizationRepository handles storage of organizations.
type OrganizationRepository struct {
	db *sql.DB
}

// NewOrganizationRepository creates a new organization repository.
func NewOrganizationRepository(db *sql.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

// Create inserts a new organization.
func (r *OrganizationRepository) Create(ctx context.Context, o *Organization) error {
	query := `
		INSERT INTO organizations (id, slug, display_name, description, avatar_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`
	row := r.db.QueryRowContext(ctx, query,
		o.ID,
		o.Slug,
		o.DisplayName,
		o.Description,
		o.AvatarURL,
		o.CreatedAt,
		o.UpdatedAt,
	)
	return row.Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt)
}

// GetBySlug retrieves an organization by its slug.
func (r *OrganizationRepository) GetBySlug(ctx context.Context, slug string) (*Organization, error) {
	query := `
		SELECT id, slug, display_name, description, avatar_url, created_at, updated_at
		FROM organizations
		WHERE slug = $1`
	row := r.db.QueryRowContext(ctx, query, slug)

	var o Organization
	err := row.Scan(
		&o.ID,
		&o.Slug,
		&o.DisplayName,
		&o.Description,
		&o.AvatarURL,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// GetByID retrieves an organization by ID.
func (r *OrganizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*Organization, error) {
	query := `
		SELECT id, slug, display_name, description, avatar_url, created_at, updated_at
		FROM organizations
		WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var o Organization
	err := row.Scan(
		&o.ID,
		&o.Slug,
		&o.DisplayName,
		&o.Description,
		&o.AvatarURL,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// List returns a list of organizations with pagination.
func (r *OrganizationRepository) List(ctx context.Context, limit, offset int) ([]*Organization, error) {
	query := `
		SELECT id, slug, display_name, description, avatar_url, created_at, updated_at
		FROM organizations
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []*Organization
	for rows.Next() {
		var o Organization
		if err := rows.Scan(
			&o.ID,
			&o.Slug,
			&o.DisplayName,
			&o.Description,
			&o.AvatarURL,
			&o.CreatedAt,
			&o.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orgs = append(orgs, &o)
	}
	return orgs, rows.Err()
}

// Update updates an organization.
func (r *OrganizationRepository) Update(ctx context.Context, o *Organization) error {
	query := `
		UPDATE organizations
		SET slug = $2, display_name = $3, description = $4, avatar_url = $5, updated_at = NOW()
		WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query,
		o.ID,
		o.Slug,
		o.DisplayName,
		o.Description,
		o.AvatarURL,
	)
	return err
}

// Delete removes an organization.
func (r *OrganizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM organizations WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}