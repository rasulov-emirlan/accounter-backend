package entities

import (
	"time"

	"github.com/google/uuid"
)

type Store struct {
	ID          uuid.UUID   `json:"id"`
	Owner       *Owner      `json:"owner,omitempty"`
	Sellers     []Seller    `json:"sellers,omitempty"`
	Warehouses  []Warehouse `json:"warehouses,omitempty"`
	Name        string      `json:"name" validate:"required"`
	Description string      `json:"description" validate:"required"`
	CreatedAt   time.Time   `json:"createdAt"`
}

func NewStore(owner *Owner, name, description string) Store {
	return Store{
		ID:          uuid.New(),
		Owner:       owner,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
	}
}
