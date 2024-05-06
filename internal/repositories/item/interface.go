package item

import (
	"context"
	"github.com/google/uuid"
	"github.com/s0vunia/password-manager/internal/domain"
)

type Repository interface {
	GetItem(ctx context.Context, userId uuid.UUID) (*domain.Item, error)
	GetItems(ctx context.Context, userId uuid.UUID) ([]*domain.Item, error)
}
