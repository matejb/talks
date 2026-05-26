package library

import "fmt"

// User represents a library member.
type User struct {
	ID    int
	Name  string
	Email string
}

type UserStore interface {
	GetUser(id int) (User, error)
	ListUsers() ([]User, error)
	CreateUser(user User) error
}

type UserService struct {
	Store UserStore
}

func (s *UserService) User(id int) (User, error) {
	return s.Store.GetUser(id)
}

func (s *UserService) Users() ([]User, error) {
	return s.Store.ListUsers()
}

func (s *UserService) CreateUser(user User) error {
	if user.Name == "" {
		return fmt.Errorf("user name is required")
	}
	return s.Store.CreateUser(user)
}
