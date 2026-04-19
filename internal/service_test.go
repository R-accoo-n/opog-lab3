package internal

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type mockProductStorage struct {
	product   Product
	id        uuid.UUID
	getErr    error
	createErr error
}

func (m mockProductStorage) Get(ctx context.Context, id uuid.UUID) (Product, error) {
	return m.product, m.getErr
}

func (m mockProductStorage) Create(ctx context.Context, params CreateProductPayload) (uuid.UUID, error) {
	return m.id, m.createErr
}

func TestCreateProduct(t *testing.T) {
	tests := []struct {
		name    string
		payload CreateProductPayload
		wantErr bool
	}{
		{
			name: "success",
			payload: CreateProductPayload{
				Name: "MacBook",
				Category: Category{
					Name: "Laptops",
					Tax:  20,
				},
				Price: 5000,
			},
			wantErr: false,
		},
		{
			name: "fail empty name",
			payload: CreateProductPayload{
				Name: "",
				Category: Category{
					Name: "Laptops",
					Tax:  20,
				},
				Price: 5000,
			},
			wantErr: true,
		},
		{
			name: "fail negative price",
			payload: CreateProductPayload{
				Name: "MacBook",
				Category: Category{
					Name: "Laptops",
					Tax:  20,
				},
				Price: -1,
			},
			wantErr: true,
		},
		{
			name: "edge zero price",
			payload: CreateProductPayload{
				Name: "Free item",
				Category: Category{
					Name: "Accessories",
					Tax:  10,
				},
				Price: 0,
			},
			wantErr: false,
		},
		{
			name: "fail negative tax",
			payload: CreateProductPayload{
				Name: "Phone",
				Category: Category{
					Name: "Smartphones",
					Tax:  -5,
				},
				Price: 1000,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewProducts(mockProductStorage{id: uuid.New()})
			_, err := svc.CreateProduct(context.Background(), tt.payload)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestGetProduct(t *testing.T) {
	id := uuid.New()

	svc := NewProducts(mockProductStorage{
		product: Product{
			ID:   id,
			Name: "MacBook",
			Category: Category{
				Name: "Laptops",
				Tax:  20,
			},
			Price: 5000,
		},
	})

	product, finalPrice, err := svc.GetProduct(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, "MacBook", product.Name)
	require.Equal(t, 6000.0, finalPrice)
}
