package app

import (
	"context"
	"github.com/s0vunia/password-manager/internal/domain"
)

type Repository interface {
	App(ctx context.Context, appID int) (domain.App, error)
}
