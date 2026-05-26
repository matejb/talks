package port

import (
	"context"

	"github.com/matejb/talks/2026/gops-2026/layouts-examples/07-hexagonal/internal/core/domain"
)

// ReviewRepository is the outbound port for review persistence.
type ReviewRepository interface {
	ByBookID(ctx context.Context, bookID int) ([]domain.Review, error)
	Save(ctx context.Context, review domain.Review) error
}
