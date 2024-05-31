package loginItem

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/s0vunia/password-manager/internal/domain"
	"github.com/s0vunia/password-manager/internal/lib/logger/sl"
	"log/slog"
)

type ILoginItemService interface {
	CreateLoginItem(ctx context.Context, item domain.LoginItem) (uuid.UUID, error)
	GetLoginItem(ctx context.Context, itemId, userId uuid.UUID) (*domain.LoginItem, error)
	GetLoginItems(ctx context.Context, userId uuid.UUID) ([]*domain.LoginItem, error)
	DeleteLoginItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error
}

type Service struct {
	log               *slog.Logger
	loginItemSaver    Saver
	loginItemProvider Provider
}

type Saver interface {
	CreateLoginItem(ctx context.Context,
		item domain.LoginItem,
	) (uuid.UUID, error)
}

type Provider interface {
	GetLoginItem(ctx context.Context,
		itemId, userId uuid.UUID,
	) (*domain.LoginItem, error)
	GetLoginItems(ctx context.Context,
		userId uuid.UUID) ([]*domain.LoginItem, error)
	DeleteLoginItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error
}

func New(
	log *slog.Logger,
	saver Saver,
	provider Provider,
) *Service {
	return &Service{
		log:               log,
		loginItemSaver:    saver,
		loginItemProvider: provider,
	}
}

func (l *Service) CreateLoginItem(ctx context.Context, item domain.LoginItem) (uuid.UUID, error) {
	const op = "LoginItemService.CreateLoginItem"

	log := l.log.With(
		slog.String("op", op),
		slog.Any("item", item),
	)

	log.Info("attempting to create login item")

	id, err := l.loginItemSaver.CreateLoginItem(ctx, item)
	if err != nil {
		log.Error("failed to create login item", sl.Err(err))

		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (l *Service) GetLoginItem(ctx context.Context, itemId, userId uuid.UUID) (*domain.LoginItem, error) {
	const op = "LoginItemService.GetLoginItem"

	log := l.log.With(
		slog.String("op", op),
		slog.String("user", userId.String()),
		slog.String("item", itemId.String()),
	)

	log.Info("attempting to get login item")
	item, err := l.loginItemProvider.GetLoginItem(ctx, itemId, userId)
	if err != nil {
		log.Error("failed to get login item", sl.Err(err))

		return &domain.LoginItem{}, fmt.Errorf("%s: %w", op, err)
	}
	return item, nil

}

func (l *Service) GetLoginItems(ctx context.Context, userId uuid.UUID) ([]*domain.LoginItem, error) {
	const op = "LoginItemService.GetLoginItems"

	log := l.log.With(
		slog.String("op", op),
		slog.String("user", userId.String()),
	)

	log.Info("attempting to get login itemS")
	items, err := l.loginItemProvider.GetLoginItems(ctx, userId)

	if err != nil {
		log.Error("failed to get login itemS", sl.Err(err))
		return make([]*domain.LoginItem, 0), fmt.Errorf("%s: %w", op, err)
	}
	return items, nil
}
func (l *Service) DeleteLoginItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error {
	const op = "LoginItemService.DeleteLoginItem"

	log := l.log.With(
		slog.String("op", op),
		slog.String("user", userId.String()),
	)

	log.Info("attempting to delete login item")
	err := l.loginItemProvider.DeleteLoginItem(ctx, userId, itemId)
	if err != nil {
		log.Error("failed to delete login item", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
