package item

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/s0vunia/password-manager/internal/domain"
	"github.com/s0vunia/password-manager/internal/lib/logger/sl"
	"log/slog"
)

type IItemService interface {
	GetItem(ctx context.Context, itemId, userId uuid.UUID) (*domain.Item, error)
	GetItems(ctx context.Context, userId uuid.UUID) ([]*domain.Item, error)
}
type Service struct {
	log          *slog.Logger
	itemProvider Provider
}

type Provider interface {
	GetItem(ctx context.Context, itemId, userId uuid.UUID) (*domain.Item, error)
	GetItems(ctx context.Context, userId uuid.UUID) ([]*domain.Item, error)
}

func New(
	log *slog.Logger,
	provider Provider,
) *Service {
	return &Service{
		log:          log,
		itemProvider: provider,
	}
}

func (s *Service) GetItem(ctx context.Context, itemId, userId uuid.UUID) (*domain.Item, error) {
	const op = "itemService.GetItem"

	log := s.log.With(
		slog.String("op", op),
		slog.String("user", userId.String()),
		slog.String("item", itemId.String()),
	)

	log.Info("attempting to get item")
	item, err := s.itemProvider.GetItem(ctx, itemId, userId)
	if err != nil {
		log.Error("failed to get login item", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return item, nil
}

func (s *Service) GetItems(ctx context.Context, userId uuid.UUID) ([]*domain.Item, error) {
	const op = "itemService.GetItems"

	log := s.log.With(
		slog.String("op", op),
		slog.String("user", userId.String()),
	)

	log.Info("attempting to get items")
	items, err := s.itemProvider.GetItems(ctx, userId)
	if err != nil {
		log.Error("failed to get login item", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return items, nil
}
