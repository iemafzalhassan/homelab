package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Session represents a user session.
type Session struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	RefreshToken string
	ExpiresAt  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// SessionRepository handles storage of sessions.
type SessionRepository struct {
	db *sql.DB
}

// NewSessionRepository creates a new session repository.
func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Create inserts a new session.
func (r *SessionRepository) Create(ctx context.Context, s *Session) error {
	query := `
		INSERT INTO sessions (id, user_id, refresh_token, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`
	row := r.db.QueryRowContext(ctx, query,
		s.ID,
		s.UserID,
		s.RefreshToken,
		s.ExpiresAt,
		s.CreatedAt,
		s.UpdatedAt,
	)
	return row.Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

// GetByRefreshToken retrieves a session by refresh token.
func (r *SessionRepository) GetByRefreshToken(ctx context.Context, token string) (*Session, error) {
	query := `
		SELECT id, user_id, refresh_token, expires_at, created_at, updated_at
		FROM sessions
		WHERE refresh_token = $1`
	row := r.db.QueryRowContext(ctx, query, token)

	var s Session
	err := row.Scan(
		&s.ID,
		&s.UserID,
		&s.RefreshToken,
		&s.ExpiresAt,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// GetByUserID retrieves sessions for a user (maybe multiple).
func (r *SessionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Session, error) {
	query := `
		SELECT id, user_id, refresh_token, expires_at, created_at, updated_at
		FROM sessions
		WHERE user_id = $1
		ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*Session
	for rows.Next() {
		var s Session
		if err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.RefreshToken,
			&s.ExpiresAt,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, &s)
	}
	return sessions, rows.Err()
}

// UpdateRefreshToken updates the refresh token and expiry for a session.
func (r *SessionRepository) UpdateRefreshToken(ctx context.Context, id uuid.UUID, token string, expiresAt time.Time) error {
	query := `
		UPDATE sessions
		SET refresh_token = $2, expires_at = $3, updated_at = NOW()
		WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id, token, expiresAt)
	return err
}

// Delete removes a session.
func (r *SessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM sessions WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// DeleteByUserID removes all sessions for a user.
func (r *SessionRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM sessions WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

// DeleteExpired removes expired sessions.
func (r *SessionRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM sessions WHERE expires_at < NOW()`
	_, err := r.db.ExecContext(ctx, query)
	return err
}