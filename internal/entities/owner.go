package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const hashCost = bcrypt.MinCost

var (
	ErrPasswordTooShort = errors.New("пароль не может содержать менее 5 символов")
)

type Owner struct {
	ID          uuid.UUID `json:"id"`
	PhoneNumber string    `json:"phoneNumber" validate:"max=500"`
	FullName    string    `json:"fullName" validate:"required"`
	Username    string    `json:"username" validate:"required,max=500"`
	Password    string    `json:"-"`
	Sellers     []Seller  `json:"sellers,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

func NewOwner(phoneNumber, fullName, username, password string) (Owner, error) {
	if len(password) < 5 {
		return Owner{}, ErrPasswordTooShort
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	if err != nil {
		return Owner{}, err
	}

	return Owner{
		ID:          uuid.New(),
		PhoneNumber: phoneNumber,
		FullName:    fullName,
		Username:    username,
		Password:    string(hashedPassword),
	}, nil
}
