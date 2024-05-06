package loginItem

import (
	"context"
	"github.com/google/uuid"
	"github.com/s0vunia/password-manager/internal/domain"
)

type Repository interface {
	CreateLoginItem(ctx context.Context, item domain.LoginItem, userId uuid.UUID) (uuid.UUID, error)
	GetLoginItem(ctx context.Context, itemId, userId uuid.UUID) (*domain.LoginItem, error)
	GetLoginItems(ctx context.Context, userId uuid.UUID) ([]*domain.LoginItem, error)
}
