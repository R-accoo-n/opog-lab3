package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/exp/slog"

	"github.com/DenisGoldiner/webapp/internal"
	"github.com/DenisGoldiner/webapp/internal/adapters/postgres"
	"github.com/DenisGoldiner/webapp/internal/ports/rest"
)

func main() {
	app := newApplication()

	slog.Info("Starting application")

	app.start()
}

type application struct {
	server *http.Server
}

func newApplication() application {
	dbExec, err := newDB()
	if err != nil {
		log.Fatal(err)
	}

	server := newServer(dbExec)

	return application{
		server: server,
	}
}

func newDB() (sqlx.ExtContext, error) {
	dsn := "postgres://postgres:postgres@localhost:5432/travellers?sslmode=disable"
	conn, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func newServer(dbExec sqlx.ExtContext) *http.Server {
	travellersClient := postgres.NewClient(dbExec)
	travellersService := internal.NewTravellers(travellersClient)

	handlers := map[string]http.Handler{
		"/api/v1/travellers": rest.NewTravellerHandler(travellersService),
	}

	mux := http.NewServeMux()
	for route, handler := range handlers {
		mux.Handle(route, handler)
	}

	return &http.Server{
		Addr:    "localhost:8081",
		Handler: mux,
	}
}

func (app application) start() {
	if err := app.server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return
		}

		log.Printf("failed to start the HTTP server, error: %v", err)
	}
}
