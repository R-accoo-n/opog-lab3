package internal

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Products struct {
	productStorage ProductStorage
}

func NewProducts(db ProductStorage) Products {
	return Products{productStorage: db}
}

func (p Products) CreateProduct(ctx context.Context, params CreateProductPayload) (uuid.UUID, error) {
	if params.Name == "" {
		return uuid.Nil, fmt.Errorf("%w: name is required", ErrInvalidInput)
	}
	if params.Category.Name == "" {
		return uuid.Nil, fmt.Errorf("%w: category name is required", ErrInvalidInput)
	}
	if params.Price < 0 {
		return uuid.Nil, fmt.Errorf("%w: price cannot be negative", ErrInvalidInput)
	}
	if params.Category.Tax < 0 {
		return uuid.Nil, fmt.Errorf("%w: tax cannot be negative", ErrInvalidInput)
	}

	id, err := p.productStorage.Create(ctx, params)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create product: %w", err)
	}

	return id, nil
}

func (p Products) GetProduct(ctx context.Context, id uuid.UUID) (Product, float64, error) {
	if id == uuid.Nil {
		return Product{}, 0, fmt.Errorf("%w: invalid uuid", ErrInvalidInput)
	}

	product, err := p.productStorage.Get(ctx, id)
	if err != nil {
		return Product{}, 0, fmt.Errorf("failed to get product: %w", err)
	}

	finalPrice := product.Price + product.Price*product.Category.Tax/100
	return product, finalPrice, nil
}
