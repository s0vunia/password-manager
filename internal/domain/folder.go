package domain

import "github.com/google/uuid"

type Folder struct {
	ID     uuid.UUID
	UserId uuid.UUID
	Name   string
}
