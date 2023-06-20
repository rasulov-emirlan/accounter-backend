package entities

import (
	"time"

	"github.com/google/uuid"
)

type Seller struct {
	ID        uuid.UUID `json:"id"`
	Owner     *Owner    `json:"owner,omitempty"`
	Username  string    `json:"username" validate:"required,max=500"`
	FullName  string    `json:"fullName" validate:"required"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewSeller(owner *Owner, username, fullName string) Seller {
	return Seller{
		ID:        uuid.New(),
		Owner:     owner,
		Username:  username,
		FullName:  fullName,
		CreatedAt: time.Now(),
	}
}
