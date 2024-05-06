package manager

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/s0vunia/password-manager/internal/domain"
	"github.com/s0vunia/password-manager/internal/lib/logger/sl"
	"log/slog"
)

type ILoginItemService interface {
	CreateLoginItem(ctx context.Context, item domain.LoginItem, userId uuid.UUID) (uuid.UUID, error)
	GetLoginItem(ctx context.Context, itemId, userId uuid.UUID) (*domain.LoginItem, error)
	GetLoginItems(ctx context.Context, userId uuid.UUID) ([]*domain.LoginItem, error)
}

type LoginItemService struct {
	log               *slog.Logger
	loginItemSaver    LoginItemSaver
	loginItemProvider LoginItemProvider
}

type LoginItemSaver interface {
	SaveItem(ctx context.Context,
		item domain.LoginItem,
		userId uuid.UUID,
	) (uuid.UUID, error)
}

type LoginItemProvider interface {
	GetItem(ctx context.Context,
		itemId, userId uuid.UUID,
	) (*domain.LoginItem, error)
	GetItems(ctx context.Context,
		userId uuid.UUID) ([]*domain.LoginItem, error)
}

func New(
	log *slog.Logger,
	saver LoginItemSaver,
	provider LoginItemProvider,
) *LoginItemService {
	return &LoginItemService{
		log:               log,
		loginItemSaver:    saver,
		loginItemProvider: provider,
	}
}

func (l *LoginItemService) CreateLoginItem(ctx context.Context, item domain.LoginItem, userId uuid.UUID) (uuid.UUID, error) {
	const op = "LoginItemService.CreateLoginItem"

	log := l.log.With(
		slog.String("op", op),
		slog.String("user", userId.String()),
		slog.Any("item", item),
	)

	log.Info("attempting to create login item")
	id, err := l.loginItemSaver.SaveItem(ctx, item, userId)
	if err != nil {
		log.Error("failed to create login item", sl.Err(err))

		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (l *LoginItemService) GetLoginItem(ctx context.Context, itemId, userId uuid.UUID) (*domain.LoginItem, error) {
	const op = "LoginItemService.GetLoginItem"

	log := l.log.With(
		slog.String("op", op),
		slog.String("user", userId.String()),
		slog.String("item", itemId.String()),
	)

	log.Info("attempting to get login item")
	item, err := l.loginItemProvider.GetItem(ctx, itemId, userId)
	if err != nil {
		log.Error("failed to get login item", sl.Err(err))

		return &domain.LoginItem{}, fmt.Errorf("%s: %w", op, err)
	}
	return item, nil

}

func (l *LoginItemService) GetLoginItems(ctx context.Context, userId uuid.UUID) ([]*domain.LoginItem, error) {
	const op = "LoginItemService.GetLoginItems"

	log := l.log.With(
		slog.String("op", op),
		slog.String("user", userId.String()),
	)

	log.Info("attempting to get login itemS")
	items, err := l.loginItemProvider.GetItems(ctx, userId)

	if err != nil {
		log.Error("failed to get login itemS", sl.Err(err))
		return make([]*domain.LoginItem, 0), fmt.Errorf("%s: %w", op, err)
	}
	return items, nil
}
