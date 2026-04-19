package main

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/R-accoo-n/opog-lab3/internal"
	"github.com/R-accoo-n/opog-lab3/internal/adapters/postgres"
	"github.com/R-accoo-n/opog-lab3/internal/ports/rest"
)

func main() {
	dsn := "postgres://postgres:postgres@localhost:5432/travellers?sslmode=disable"

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	server := newServer(db)
	fmt.Println("Server started on http://localhost:8081")

	if err = server.ListenAndServe(); err != nil {
		panic(fmt.Errorf("server stopped: %w", err))
	}
}

func newServer(db *sqlx.DB) *http.Server {
	productsClient := postgres.NewClient(db)
	productsService := internal.NewProducts(productsClient)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/products", rest.NewProductHandler(productsService))

	return &http.Server{
		Addr:    "localhost:8081",
		Handler: mux,
	}
}
