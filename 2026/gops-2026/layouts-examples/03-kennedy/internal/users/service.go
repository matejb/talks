package users

import (
	"fmt"
	"sync"
)

// Service provides user business logic.
// In the Kennedy layout, domain services directly hold a database
// connection - no interface abstraction layer.
type Service struct {
	mu     sync.RWMutex
	users  map[int]User
	nextID int
}

// NewService creates a user service with pre-seeded data.
func NewService() *Service {
	return &Service{
		users: map[int]User{
			1: {ID: 1, Name: "Alice", Email: "alice@example.com"},
			2: {ID: 2, Name: "Bob", Email: "bob@example.com"},
		},
		nextID: 3,
	}
}

// User returns a user by ID.
func (s *Service) User(id int) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	if !ok {
		return User{}, fmt.Errorf("user %d not found", id)
	}
	return u, nil
}

// Users returns all users.
func (s *Service) Users() ([]User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]User, 0, len(s.users))
	for _, u := range s.users {
		result = append(result, u)
	}
	return result, nil
}

// CreateUser adds a new user.
func (s *Service) CreateUser(u User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	u.ID = s.nextID
	s.nextID++
	s.users[u.ID] = u
	return nil
}
