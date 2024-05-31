package user

import (
	"context"
	"github.com/google/uuid"
	"github.com/s0vunia/password-manager/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, login string, passHash []byte) (uuid.UUID, error)
	Get(ctx context.Context, login string) (domain.User, error)
}
