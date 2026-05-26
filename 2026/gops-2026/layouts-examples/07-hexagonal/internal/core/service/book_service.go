package service

import (
	"context"
	"fmt"

	"github.com/matejb/talks/2026/gops-2026/layouts-examples/07-hexagonal/internal/core/domain"
	"github.com/matejb/talks/2026/gops-2026/layouts-examples/07-hexagonal/internal/core/port"
)

// BookService implements the book-related use cases (inbound port).
type BookService struct {
	books port.BookRepository
	users port.UserRepository
}

// NewBookService creates a new BookService with the given repositories.
func NewBookService(books port.BookRepository, users port.UserRepository) *BookService {
	return &BookService{books: books, users: users}
}

// Get returns a single book by ID.
func (s *BookService) Get(ctx context.Context, id int) (domain.Book, error) {
	return s.books.ByID(ctx, id)
}

// All returns all books.
func (s *BookService) All(ctx context.Context) ([]domain.Book, error) {
	return s.books.All(ctx)
}

// Borrow checks out a book to a user.
func (s *BookService) Borrow(ctx context.Context, bookID, userID int) error {
	book, err := s.books.ByID(ctx, bookID)
	if err != nil {
		return fmt.Errorf("find book: %w", err)
	}
	if book.Borrower != nil {
		return fmt.Errorf("book %d already borrowed by %s", bookID, book.Borrower.Name)
	}

	user, err := s.users.ByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("find user: %w", err)
	}

	book.Borrower = &domain.User{ID: user.ID, Name: user.Name, Email: user.Email}
	return s.books.Update(ctx, book)
}

// Return returns a borrowed book.
func (s *BookService) Return(ctx context.Context, bookID int) error {
	book, err := s.books.ByID(ctx, bookID)
	if err != nil {
		return fmt.Errorf("find book: %w", err)
	}
	if book.Borrower == nil {
		return fmt.Errorf("book %d is not borrowed", bookID)
	}
	book.Borrower = nil
	return s.books.Update(ctx, book)
}
