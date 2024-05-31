package app

import (
	"context"
	"github.com/s0vunia/password-manager/internal/domain"
)

type Repository interface {
	App(ctx context.Context, appID int64) (domain.App, error)
}
