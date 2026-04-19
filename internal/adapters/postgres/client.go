package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/R-accoo-n/opog-lab3/internal"
)

type Client struct {
	dbExec sqlx.ExtContext
}

func NewClient(dbExec sqlx.ExtContext) Client {
	return Client{dbExec: dbExec}
}

func (c Client) Get(ctx context.Context, id uuid.UUID) (internal.Product, error) {
	q := `select id, name, category_name, category_tax, price
	      from products
	      where id = $1`

	rows, err := c.dbExec.QueryxContext(ctx, q, id)
	if err != nil {
		return internal.Product{}, fmt.Errorf("failed to fetch product: %w", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err = rows.StructScan(&product); err != nil {
			return internal.Product{}, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	if len(products) == 0 {
		return internal.Product{}, fmt.Errorf("no product with id %s: %w", id, internal.ErrNoResource)
	}

	return internal.Product{
		ID:   products[0].ID,
		Name: products[0].Name,
		Category: internal.Category{
			Name: products[0].CategoryName,
			Tax:  products[0].CategoryTax,
		},
		Price: products[0].Price,
	}, nil
}

func (c Client) Create(ctx context.Context, params internal.CreateProductPayload) (uuid.UUID, error) {
	q := `insert into products (name, category_name, category_tax, price)
	      values ($1, $2, $3, $4)
	      returning id`

	rows, err := c.dbExec.QueryxContext(ctx, q,
		params.Name,
		params.Category.Name,
		params.Category.Tax,
		params.Price,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create product: %w", err)
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err = rows.Scan(&id); err != nil {
			return uuid.Nil, fmt.Errorf("failed to scan product id: %w", err)
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return uuid.Nil, fmt.Errorf("failed to create product: no id returned")
	}

	return ids[0], nil
}
