package loginItem

import (
	"context"
	"github.com/google/uuid"
	"github.com/s0vunia/password-manager/internal/domain"
)

type Repository interface {
	CreateLoginItem(ctx context.Context, item domain.LoginItem) (uuid.UUID, error)
	GetLoginItem(ctx context.Context, itemId, userId uuid.UUID) (*domain.LoginItem, error)
	GetLoginItems(ctx context.Context, userId uuid.UUID) ([]*domain.LoginItem, error)
	DeleteLoginItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error
}
