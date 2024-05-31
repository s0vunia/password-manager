package repositories

import "errors"

var (
	ErrUserExists     = errors.New("user already exists")
	ErrUserNotFound   = errors.New("user not found")
	ErrAppNotFound    = errors.New("app not found")
	ErrItemNotFound   = errors.New("item not found")
	ErrFolderExists   = errors.New("folder not exists")
	ErrFolderNotFound = errors.New("folder not exists")
	ErrItemExists     = errors.New("item already exists")
)
