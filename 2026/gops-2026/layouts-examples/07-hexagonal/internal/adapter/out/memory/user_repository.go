package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/matejb/talks/2026/gops-2026/layouts-examples/07-hexagonal/internal/core/domain"
	"github.com/matejb/talks/2026/gops-2026/layouts-examples/07-hexagonal/internal/core/port"
)

// Compile-time check: UserRepository implements port.UserRepository.
var _ port.UserRepository = (*UserRepository)(nil)

// UserRepository is an in-memory implementation of port.UserRepository.
type UserRepository struct {
	mu     sync.RWMutex
	users  map[int]domain.User
	nextID int
}

// NewUserRepository creates a UserRepository with pre-seeded data.
func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: map[int]domain.User{
			1: {ID: 1, Name: "Alice", Email: "alice@example.com"},
			2: {ID: 2, Name: "Bob", Email: "bob@example.com"},
		},
		nextID: 3,
	}
}

func (r *UserRepository) ByID(_ context.Context, id int) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return domain.User{}, fmt.Errorf("user %d not found", id)
	}
	return u, nil
}

func (r *UserRepository) All(_ context.Context) ([]domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]domain.User, 0, len(r.users))
	for _, u := range r.users {
		result = append(result, u)
	}
	return result, nil
}

func (r *UserRepository) Save(_ context.Context, u domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	u.ID = r.nextID
	r.nextID++
	r.users[u.ID] = u
	return nil
}
