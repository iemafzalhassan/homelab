package service

import (
	"context"
	"errors"

	"github.com/kyros-platform/kyros/services/auth/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles user-related operations.
type UserService struct {
	repo  *repository.UserRepository
	logger *repository.Logger // we'll use a simple logger for now; later we can integrate with kyros/internal/logger
}

// NewUserService creates a new user service.
func NewUserService(repo *repository.UserRepository, logger *repository.Logger) *UserService {
	return &UserService{
		repo:  repo,
		logger: logger,
	}
}

// Register creates a new user with the given email, display name, and password.
func (s *UserService) Register(ctx context.Context, email, displayName, password string) (*repository.User, error) {
	// Check if user already exists
	existing, err := s.repo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &repository.User{
		ID:           uuid.New(),
		Email:        email,
		DisplayName:  displayName,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Authenticate checks the email and password and returns the user if valid.
func (s *UserService) Authenticate(ctx context.Context, email, password string) (*repository.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// If the user has no password hash (e.g., Keycloak user), we cannot authenticate via password.
	if user.PasswordHash == "" {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// UpdateProfile updates the user's display name and avatar URL.
func (s *UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, displayName string, avatarURL string) (*repository.User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.DisplayName = displayName
	user.AvatarURL = sql.NullString{String: avatarURL, Valid: avatarURL != ""}
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// ChangePassword changes the user's password.
func (s *UserService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Check old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("invalid old password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.repo.UpdatePassword(ctx, userID, string(hashedPassword)); err != nil {
		return err
	}
	return nil
}