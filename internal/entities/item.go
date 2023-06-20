package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSizeExclusive = errors.New("размер должен быть либо числовым диапазоном, либо символом")
)

type (
	Item struct {
		ID          uuid.UUID `json:"id"`
		Store       *Store    `json:"store,omitempty"`
		Category    *Category `json:"category,omitempty"`
		Name        string    `json:"name" validate:"required"`
		Article     string    `json:"article" validate:"required,max=100"`
		Description string    `json:"description"`
		IconURL     string    `json:"iconURL"`
		Color       string    `json:"color" validate:"required,max=6"`
		Price       float64   `json:"price" validate:"required"` // retail price
		Sizes       []Size    `json:"sizes,omitempty"`
		CreatedAt   time.Time `json:"createdAt"`
	}

	Size struct {
		ID        int64      `json:"id"`
		Item      *Item      `json:"item,omitempty"`
		Warehouse *Warehouse `json:"warehouse,omitempty"`

		// SizeNumber and SizeSymbol are mutually exclusive

		SizeNumber *string `json:"sizeNumber,omitempty"` // number range like 36-40
		SizeSymbol *string `json:"sizeSymbol,omitempty"` // symbols like S, M, L, XL

		Quantity  int64     `json:"quantity" validate:"required"`
		Cost      float64   `json:"cost" validate:"required"` // cost of all items with this size
		CreatedAt time.Time `json:"createdAt"`
	}
)

func NewItem(
	store *Store,
	category *Category,
	name,
	article,
	description,
	iconURL,
	color string,
	price float64,
) Item {

	return Item{
		ID:          uuid.New(),
		Store:       store,
		Category:    category,
		Name:        name,
		Article:     article,
		Description: description,
		IconURL:     iconURL,
		Color:       color,
		Price:       price,
		CreatedAt:   time.Now(),
	}
}

func NewSize(
	item *Item,
	warehouse *Warehouse,
	sizeNumber,
	sizeSymbol *string,
	quantity int64,
	cost float64,
) (Size, error) {
	if sizeNumber != nil && sizeSymbol != nil {
		return Size{}, ErrSizeExclusive
	}
	return Size{
		Item:       item,
		Warehouse:  warehouse,
		SizeNumber: sizeNumber,
		SizeSymbol: sizeSymbol,
		Quantity:   quantity,
		Cost:       cost,
		CreatedAt:  time.Now(),
	}, nil
}
