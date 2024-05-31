package folder

import (
	"context"
	"github.com/google/uuid"
	"github.com/s0vunia/password-manager/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, folder domain.Folder) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Folder, error)
}
