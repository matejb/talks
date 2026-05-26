package memory

import (
	"fmt"
	"sync"

	"github.com/matejb/talks/2026/gops-2026/layouts-examples/04-bourgon-pkg/pkg/library"
)

// UserService is an in-memory implementation of library.UserService.
type UserService struct {
	mu     sync.RWMutex
	users  map[int]library.User
	nextID int
}

// NewUserService returns a UserService pre-seeded with sample data.
func NewUserService() *UserService {
	return &UserService{
		users: map[int]library.User{
			1: {ID: 1, Name: "Alice", Email: "alice@example.com"},
			2: {ID: 2, Name: "Bob", Email: "bob@example.com"},
		},
		nextID: 3,
	}
}

func (s *UserService) User(id int) (library.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	if !ok {
		return library.User{}, fmt.Errorf("user %d not found", id)
	}
	return u, nil
}

func (s *UserService) Users() ([]library.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]library.User, 0, len(s.users))
	for _, u := range s.users {
		result = append(result, u)
	}
	return result, nil
}

func (s *UserService) CreateUser(user library.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	user.ID = s.nextID
	s.nextID++
	s.users[user.ID] = user
	return nil
}
