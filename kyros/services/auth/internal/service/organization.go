package service

import (
	"context"

	"github.com/kyros-platform/kyros/services/auth/internal/repository"
)

// OrganizationService handles organization-related operations.
type OrganizationService struct {
	repo  *repository.OrganizationRepository
	logger *repository.Logger
}

// NewOrganizationService creates a new organization service.
func NewOrganizationService(repo *repository.OrganizationRepository, logger *repository.Logger) *OrganizationService {
	return &OrganizationService{
		repo:  repo,
		logger: logger,
	}
}

// Create creates a new organization.
func (s *OrganizationService) Create(ctx context.Context, slug, displayName, description, avatarURL string) (*repository.Organization, error) {
	org := &repository.Organization{
		ID:          uuid.New(),
		Slug:        slug,
		DisplayName: displayName,
		Description: sql.NullString{String: description, Valid: description != ""},
		AvatarURL:   sql.NullString{String: avatarURL, Valid: avatarURL != ""},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, org); err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}
	return org, nil
}

// GetBySlug retrieves an organization by its slug.
func (s *OrganizationService) GetBySlug(ctx context.Context, slug string) (*repository.Organization, error) {
	return s.repo.GetBySlug(ctx, slug)
}

// GetByID retrieves an organization by ID.
func (s *OrganizationService) GetByID(ctx context.Context, id uuid.UUID) (*repository.Organization, error) {
	return s.repo.GetByID(ctx, id)
}

// List returns a list of organizations with pagination.
func (s *OrganizationService) List(ctx context.Context, limit, offset int) ([]*repository.Organization, error) {
	return s.repo.List(ctx, limit, offset)
}

// Update updates an organization.
func (s *OrganizationService) Update(ctx context.Context, id uuid.UUID, slug, displayName, description, avatarURL string) (*repository.Organization, error) {
	org, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	org.Slug = slug
	org.DisplayName = displayName
	org.Description = sql.NullString{String: description, Valid: description != ""}
	org.AvatarURL = sql.NullString{String: avatarURL, Valid: avatarURL != ""}
	org.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, org); err != nil {
		return nil, err
	}
	return org, nil
}

// Delete removes an organization.
func (s *OrganizationService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}