package noteItem

import (
	"context"
	"github.com/google/uuid"
	"github.com/s0vunia/password-manager/internal/domain"
	"github.com/s0vunia/password-manager/internal/services/manager/item"
)

type INoteItemService interface {
	item.IItemService
	CreateNoteItem(ctx context.Context, item *domain.NoteItem, userId uuid.UUID) (*domain.NoteItem, error)
	GetNoteItem(ctx context.Context, itemId, userId uuid.UUID) (*domain.NoteItem, error)
	GetNoteItems(ctx context.Context, userId uuid.UUID) ([]*domain.NoteItem, error)
}
