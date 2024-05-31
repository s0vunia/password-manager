package item

import (
	"context"
	"github.com/google/uuid"
	"github.com/s0vunia/password-manager/internal/domain"
)

type Repository interface {
	GetItem(ctx context.Context, itemId, userId uuid.UUID) (*domain.Item, error)
	GetItems(ctx context.Context, userId uuid.UUID) ([]*domain.Item, error)
	CreateItem(ctx context.Context, item domain.Item) (uuid.UUID, error)
}
