package domain

import "github.com/google/uuid"

type ItemType int64

const (
	ItemTypeLogin ItemType = iota
	ItemTypeNote  ItemType = iota
)

type Item struct {
	ID         uuid.UUID
	Type       ItemType
	Name       string
	FolderId   uuid.UUID
	UserId     uuid.UUID
	IsFavorite bool
}
