package categories

import "github.com/rasulov-emirlan/esep-backend/internal/entities"

type (
	CreateInput struct {
		Name             string  `json:"name" validate:"required,max=255"`
		Article          *string `json:"article" validate:"max=100"`
		StoreID          string  `json:"storeID" validate:"required,uuid4"`
		ParentCategoryID *string `json:"parentCategoryID,omitempty" validate:"uuid4"`
	}

	ReadByInput struct {
		// if ID is set, other filters will be ignored
		ID               entities.OptField[string] `json:"id" validate:"uuid4"`
		StoreID          entities.OptField[string] `json:"storeID" validate:"uuid4"`
		Text             entities.OptField[string] `json:"text" validate:"max=255"`
		ParentCategoryID entities.OptField[string] `json:"parentCategoryID" validate:"uuid4"`

		// Pagination
		PageNumber entities.OptField[uint64] `json:"pageNumber"`
		PageSize   entities.OptField[uint]   `json:"pageSize"`

		// Sorting
		SortBy    entities.OptField[string] `json:"sortBy"`    // name, article, createdAt
		SortOrder entities.OptField[string] `json:"sortOrder"` // asc, desc
	}

	UpdateInput struct {
		ID               string                     `json:"id" validate:"required,uuid4"`
		Name             entities.OptField[string]  `json:"name" validate:"max=255"`
		Article          entities.OptField[*string] `json:"article" validate:"max=100"`
		ParentCategoryID entities.OptField[*string] `json:"parentCategoryID,omitempty" validate:"uuid4"`
	}
)
