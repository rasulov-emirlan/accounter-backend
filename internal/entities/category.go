package entities

import (
	"time"

	"github.com/google/uuid"
)

type (
	DefaultCategory struct {
		ID        int64     `json:"id"`
		Name      string    `json:"name"`
		IconURL   string    `json:"iconURL"`
		CreatedAt time.Time `json:"createdAt"`
	}

	Category struct {
		ID             uuid.UUID `json:"id"`
		Store          *Store    `json:"store,omitempty"`
		ParentCategory *Category `json:"parentCategory,omitempty"`
		Name           string    `json:"name" validate:"required,max=255"`
		Article        *string   `json:"article" validate:"required,max=100"`
		IconURL        string    `json:"iconURL"`
		CreatedAt      time.Time `json:"createdAt"`
	}
)
