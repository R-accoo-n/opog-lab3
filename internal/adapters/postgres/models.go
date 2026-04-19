package postgres

import "github.com/google/uuid"

type Product struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	CategoryName string    `db:"category_name"`
	CategoryTax  float64   `db:"category_tax"`
	Price        float64   `db:"price"`
}
