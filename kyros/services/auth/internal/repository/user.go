package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system.
type User struct {
	ID           uuid.UUID
	KeycloakSub  sql.NullString // nullable after migration
	Email        string
	DisplayName  string
	AvatarURL    sql.NullString
	PasswordHash string // hashed password, empty if using external auth (like Keycloak)
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// UserRepository handles storage of users.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user into the database.
func (r *UserRepository) Create(ctx context.Context, u *User) error {
	query := `
		INSERT INTO users (id, keycloak_sub, email, display_name, avatar_url, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`
	row := r.db.QueryRowContext(ctx, query,
		u.ID,
		u.KeycloakSub,
		u.Email,
		u.DisplayName,
		u.AvatarURL,
		u.PasswordHash,
		u.CreatedAt,
		u.UpdatedAt,
	)
	return row.Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

// GetByEmail retrieves a user by email.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, keycloak_sub, email, display_name, avatar_url, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, email)

	var u User
	err := row.Scan(
		&u.ID,
		&u.KeycloakSub,
		&u.Email,
		&u.DisplayName,
		&u.AvatarURL,
		&u.PasswordHash,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetByID retrieves a user by ID.
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	query := `
		SELECT id, keycloak_sub, email, display_name, avatar_url, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var u User
	err := row.Scan(
		&u.ID,
		&u.KeycloakSub,
		&u.Email,
		&u.DisplayName,
		&u.AvatarURL,
		&u.PasswordHash,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// UpdatePassword updates the password hash for a user.
func (r *UserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, hash string) error {
	query := `
		UPDATE users
		SET password_hash = $1, updated_at = NOW()
		WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, hash, id)
	return err
}

// Update updates a user's profile information.
func (r *UserRepository) Update(ctx context.Context, u *User) error {
	query := `
		UPDATE users
		SET keycloak_sub = $2, email = $3, display_name = $4, avatar_url = $5, password_hash = $6, updated_at = NOW()
		WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query,
		u.ID,
		u.KeycloakSub,
		u.Email,
		&u.DisplayName,
		&u.AvatarURL,
		u.PasswordHash,
	)
	return err
}

// Delete removes a user from the database.
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List returns a list of users with pagination.
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*User, error) {
	query := `
		SELECT id, keycloak_sub, email, display_name, avatar_url, password_hash, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var u User
		if err := rows.Scan(
			&u.ID,
			&u.KeycloakSub,
			&u.Email,
			&u.DisplayName,
			&u.AvatarURL,
			&u.PasswordHash,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, rows.Err()
}