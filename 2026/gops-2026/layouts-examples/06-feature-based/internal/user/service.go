package user

import "context"

// Service provides user business logic.
type Service struct {
	repo Repository
}

// NewService creates a new user Service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// User returns a single user by ID.
func (s *Service) User(_ context.Context, id int) (User, error) {
	return s.repo.ByID(id)
}

// Users returns all users.
func (s *Service) Users(_ context.Context) ([]User, error) {
	return s.repo.All()
}

// CreateUser adds a new user.
func (s *Service) CreateUser(_ context.Context, u User) error {
	return s.repo.Create(u)
}
