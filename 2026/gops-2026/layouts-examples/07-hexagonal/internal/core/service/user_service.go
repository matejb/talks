package service

import (
	"context"

	"github.com/matejb/talks/2026/gops-2026/layouts-examples/07-hexagonal/internal/core/domain"
	"github.com/matejb/talks/2026/gops-2026/layouts-examples/07-hexagonal/internal/core/port"
)

// UserService implements the user-related use cases.
type UserService struct {
	users port.UserRepository
}

// NewUserService creates a new UserService with the given repository.
func NewUserService(users port.UserRepository) *UserService {
	return &UserService{users: users}
}

// Get returns a single user by ID.
func (s *UserService) Get(ctx context.Context, id int) (domain.User, error) {
	return s.users.ByID(ctx, id)
}

// All returns all users.
func (s *UserService) All(ctx context.Context) ([]domain.User, error) {
	return s.users.All(ctx)
}

// Create adds a new user.
func (s *UserService) Create(ctx context.Context, u domain.User) error {
	return s.users.Save(ctx, u)
}
