package app

import (
	grpcapp "github.com/s0vunia/password-manager/internal/app/grpc"
	"github.com/s0vunia/password-manager/internal/repositories/app"
	"github.com/s0vunia/password-manager/internal/services/auth"
	"github.com/s0vunia/password-manager/internal/services/manager/item"
	"github.com/s0vunia/password-manager/internal/services/manager/loginItem"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	item item.IItemService,
	loginItem loginItem.ILoginItemService,
	appRepo app.Repository,
	auth auth.IOAuth,
	grpcPort int,
) *App {
	grpcServer := grpcapp.New(log, auth, item, loginItem, appRepo, grpcPort)
	return &App{
		GRPCServer: grpcServer,
	}
}
