package port

import (
	"context"

	"github.com/farzadamr/go-backend-api/internal/domain"
)

type ItemRepository interface {
	Create(ctx context.Context, item *domain.Item) (*domain.Item, error)
	GetByID(ctx context.Context, id int64) (*domain.Item, error)
	Update(ctx context.Context, item *domain.Item) (*domain.Item, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, filter ItemFilter) ([]*domain.Item, int64, error)
}

type ItemFilter struct {
	CategoryId  int64
	IsAvailable *bool
}
