package port

import (
	"context"

	"github.com/matejb/talks/2026/gops-2026/layouts-examples/07-hexagonal/internal/core/domain"
)

// UserRepository is the outbound port for user persistence.
type UserRepository interface {
	ByID(ctx context.Context, id int) (domain.User, error)
	All(ctx context.Context) ([]domain.User, error)
	Save(ctx context.Context, user domain.User) error
}
