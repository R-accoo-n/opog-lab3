package internal

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrNoResource = errors.New("no resource found")
var ErrAlreadyExists = errors.New("resource already exists")
var ErrInvalidInput = errors.New("invalid input")

type ProductStorage interface {
	Get(ctx context.Context, id uuid.UUID) (Product, error)
	Create(ctx context.Context, params CreateProductPayload) (uuid.UUID, error)
}

type Category struct {
	Name string
	Tax  float64
}

type Product struct {
	ID       uuid.UUID
	Name     string
	Category Category
	Price    float64
}

type CreateProductPayload struct {
	Name     string
	Category Category
	Price    float64
}
