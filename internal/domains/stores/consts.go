package stores

import "errors"

const (
	SortByCreatedAt = "createdAt"
	SortByName      = "name"
)

var ErrDefault = errors.New("что-то пошло не так")
