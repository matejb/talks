// Package library contains domain types and service interfaces.
// It has zero external dependencies.
package library

// User represents a library member.
type User struct {
	ID    int
	Name  string
	Email string
}

// UserService provides operations on users.
type UserService interface {
	User(id int) (User, error)
	Users() ([]User, error)
	CreateUser(user User) error
}
