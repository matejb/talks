package port

import (
	"context"

	"github.com/matejb/talks/2026/gops-2026/layouts-examples/07-hexagonal/internal/core/domain"
)

// BookRepository is the outbound port for book persistence.
type BookRepository interface {
	ByID(ctx context.Context, id int) (domain.Book, error)
	All(ctx context.Context) ([]domain.Book, error)
	Save(ctx context.Context, book domain.Book) error
	Update(ctx context.Context, book domain.Book) error
}
