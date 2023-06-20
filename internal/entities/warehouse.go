package entities

import (
	"time"

	"github.com/google/uuid"
)

type Warehouse struct {
	ID          uuid.UUID `json:"id"`
	Owner       *Owner    `json:"owner,omitempty"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"required"`
	CreatedAt   time.Time `json:"createdAt"`
}

func NewWarehouse(owner *Owner, name, description string) Warehouse {
	return Warehouse{
		ID:          uuid.New(),
		Owner:       owner,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
	}
}
