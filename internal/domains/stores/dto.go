package stores

import "github.com/rasulov-emirlan/esep-backend/internal/entities"

type (
	CreateInput struct {
		OwnerID     string `json:"ownerID" validate:"required"`
		Name        string `json:"name" validate:"required,min=3"`
		Description string `json:"description"`
	}

	ReadByInput struct {
		Text    entities.OptField[string] `json:"text"`
		OwnerID entities.OptField[string] `json:"ownerID"`

		// Pagination
		PageNumber entities.OptField[int] `json:"pageNumber"`
		PageSize   entities.OptField[int] `json:"pageSize"`

		// Sorting
		SortBy    entities.OptField[string] `json:"sortBy"`    // name, createdAt
		SortOrder entities.OptField[string] `json:"sortOrder"` // asc, desc
	}

	UpdateInput struct {
		Name        entities.OptField[string] `json:"name"`
		Description entities.OptField[string] `json:"description"`
	}
)
