package employee

import (
	"context"
)

// In a larger app, this service would handle employee-related business logic,
// like updating profiles, password resets, etc. For this project, it's
// minimal as the repository is sufficient for the auth service's needs.
type Service interface {
	// Placeholder for future methods
	Create(ctx context.Context, employee *Employee) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Create(ctx context.Context, employee *Employee) error {
	return s.repo.Create(ctx, employee)
}
