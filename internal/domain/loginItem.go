package domain

import "github.com/google/uuid"

type LoginItem struct {
	Item
	ID              uuid.UUID
	Login           string
	EncryptPassword string
}

func LoginItemToItem(model *LoginItem) *Item {
	return &Item{
		ID:         model.Item.ID,
		Type:       model.Type,
		Name:       model.Name,
		FolderId:   model.FolderId,
		UserId:     model.UserId,
		IsFavorite: model.IsFavorite,
	}
}
