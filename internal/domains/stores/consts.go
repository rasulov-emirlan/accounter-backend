package stores

import "errors"

const (
	PackageName = "internal/domains/stores/"

	SortByCreatedAt = "createdAt"
	SortByName      = "name"
)

var ErrDefault = errors.New("что-то пошло не так")
