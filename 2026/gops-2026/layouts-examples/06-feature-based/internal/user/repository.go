package user

import (
	"fmt"
	"sync"
)

// Repository defines the data access interface for users.
type Repository interface {
	ByID(id int) (User, error)
	All() ([]User, error)
	Create(u User) error
}

// In-memory implementation of Repository.
type memoryRepository struct {
	mu     sync.RWMutex
	users  map[int]User
	nextID int
}

// NewMemoryRepository returns a Repository backed by memory, pre-seeded with sample data.
func NewMemoryRepository() Repository {
	return &memoryRepository{
		users: map[int]User{
			1: {ID: 1, Name: "Alice", Email: "alice@example.com"},
			2: {ID: 2, Name: "Bob", Email: "bob@example.com"},
		},
		nextID: 3,
	}
}

func (r *memoryRepository) ByID(id int) (User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return User{}, fmt.Errorf("user %d not found", id)
	}
	return u, nil
}

func (r *memoryRepository) All() ([]User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]User, 0, len(r.users))
	for _, u := range r.users {
		result = append(result, u)
	}
	return result, nil
}

func (r *memoryRepository) Create(u User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	u.ID = r.nextID
	r.nextID++
	r.users[u.ID] = u
	return nil
}
