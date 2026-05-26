package main

import (
	"fmt"
	"sync"
)

type userService struct {
	mu     sync.RWMutex
	users  map[int]User
	nextID int
}

func newUserService() *userService {
	return &userService{
		users: map[int]User{
			1: {ID: 1, Name: "Alice", Email: "alice@example.com"},
			2: {ID: 2, Name: "Bob", Email: "bob@example.com"},
		},
		nextID: 3,
	}
}

func (s *userService) User(id int) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	if !ok {
		return User{}, fmt.Errorf("user %d not found", id)
	}
	return u, nil
}

func (s *userService) Users() ([]User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]User, 0, len(s.users))
	for _, u := range s.users {
		result = append(result, u)
	}
	return result, nil
}

func (s *userService) CreateUser(user User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	user.ID = s.nextID
	s.nextID++
	s.users[user.ID] = user
	return nil
}
