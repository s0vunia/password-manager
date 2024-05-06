package domain

type LoginItem struct {
	Item
	Login           string
	EncryptPassword string
}

func LoginItemToItem(model *LoginItem) *Item {
	return &Item{
		ID:         model.ID,
		Type:       model.Type,
		Name:       model.Name,
		FolderId:   model.FolderId,
		UserId:     model.UserId,
		IsFavorite: model.IsFavorite,
	}
}
