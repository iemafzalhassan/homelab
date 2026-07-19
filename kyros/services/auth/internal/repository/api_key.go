package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// APIKey represents an API key.
type APIKey struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	KeyHash   string // bcrypt hash of the key
	Name      string
	ExpiresAt sql.NullTime
	CreatedAt time.Time
	UpdatedAt time.Time
}

// APIKeyRepository handles storage of API keys.
type APIKeyRepository struct {
	db *sql.DB
}

// NewAPIKeyRepository creates a new API key repository.
func NewAPIKeyRepository(db *sql.DB) *APIKeyRepository {
	return &APIKeyRepository{db: db}
}

// Create inserts a new API key.
func (r *APIKeyRepository) Create(ctx context.Context, k *APIKey) error {
	query := `
		INSERT INTO api_keys (id, user_id, key_hash, name, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`
	row := r.db.QueryRowContext(ctx, query,
		k.ID,
		k.UserID,
		k.KeyHash,
		k.Name,
		k.ExpiresAt,
		k.CreatedAt,
		k.UpdatedAt,
	)
	return row.Scan(&k.ID, &k.CreatedAt, &k.UpdatedAt)
}

// GetByID retrieves an API key by ID.
func (r *APIKeyRepository) GetByID(ctx context.Context, id uuid.UUID) (*APIKey, error) {
	query := `
		SELECT id, user_id, key_hash, name, expires_at, created_at, updated_at
		FROM api_keys
		WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var k APIKey
	err := row.Scan(
		&k.ID,
		&k.UserID,
		&k.KeyHash,
		&k.Name,
		&k.ExpiresAt,
		&k.CreatedAt,
		&k.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &k, nil
}

// ListByUserID returns all API keys for a user.
func (r *APIKeyRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*APIKey, error) {
	query := `
		SELECT id, user_id, key_hash, name, expires_at, created_at, updated_at
		FROM api_keys
		WHERE user_id = $1
		ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []*APIKey
	for rows.Next() {
		var k APIKey
		if err := rows.Scan(
			&k.ID,
			&k.UserID,
			&k.KeyHash,
			&k.Name,
			&k.ExpiresAt,
			&k.CreatedAt,
			&k.UpdatedAt,
		); err != nil {
			return nil, err
		}
		keys = append(keys, &k)
	}
	return keys, rows.Err()
}

// Delete removes an API key by ID.
func (r *APIKeyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM api_keys WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// ValidateKey checks the provided plaintext key against the stored hash and returns the associated API key if valid.
func (r *APIKeyRepository) ValidateKey(ctx context.Context, key string) (*APIKey, error) {
	// We need to find the API key by comparing the hash.
	// Since we cannot reverse the hash, we need to retrieve potential candidates.
	// We'll retrieve all API keys (not efficient) but for simplicity we'll do it.
	// In a production system, we would use a lookup table (e.g., hash the key with a fast hash for lookup).
	// Given the scope, we'll assume the number of API keys is manageable.

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, key_hash, name, expires_at, created_at, updated_at
		FROM api_keys`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var k APIKey
		if err := rows.Scan(
			&k.ID,
			&k.UserID,
			&k.KeyHash,
			&k.Name,
			&k.ExpiresAt,
			&k.CreatedAt,
			&k.UpdatedAt,
		); err != nil {
			return nil, err
		}

		// Compare the provided key with the stored hash
		if err := bcrypt.CompareHashAndPassword([]byte(k.KeyHash), []byte(key)); err == nil {
			// Check expiration
			if k.ExpiresAt.Valid && time.Now().After(k.ExpiresAt.Time) {
				continue // expired
			}
			return &k, nil
		}
	}

	return nil, sql.ErrNoRows
}