package categories

import "errors"

const (
	PackageName = "internal/domains/categories/"

	// Sorting
	SortByCreatedAt = "createdAt"
	SortByArticle   = "article"
	SortByName      = "name"

	// Sorting order
	SortOrderAsc  = "asc"
	SortOrderDesc = "desc"
)

var (
	ErrDefault = errors.New("что-то пошло не так")
)
