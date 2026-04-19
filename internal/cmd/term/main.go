package main

import (
	"bytes"
	"context"
	"log"
	"os"
	"runtime/pprof"
	"runtime/trace"

	"github.com/DenisGoldiner/webapp/internal"
	"github.com/DenisGoldiner/webapp/internal/adapters/postgres"
	"github.com/DenisGoldiner/webapp/internal/ports/ftp"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// go tool pprof -http=:8080 cpu_profile.prof
// go tool trace trace.out

func main() {
	traceRun()
}

func simpleRun() {
	dbExec, err := newDB()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	travellersClient := postgres.NewClient(dbExec)
	travellersService := internal.NewTravellers(travellersClient)
	travellersParser := ftp.NewParser(travellersService)

	if err = travellersParser.Run(ctx, "/Users/denys/Go/src/github.com/DenisGoldiner/webapp/internal/integration/data/test_1.csv"); err != nil {
		log.Printf("error running travellers import: %v", err)
		return
	}
}

func concurrentRun() {
	dbExec, err := newDB()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	travellersClient := postgres.NewClient(dbExec)
	travellersService := internal.NewTravellers(travellersClient)
	travellersParser := ftp.NewConcurrentParser(travellersService)

	if err = travellersParser.Run(ctx, "/Users/denys/Go/src/github.com/DenisGoldiner/webapp/internal/integration/data/test_1.csv"); err != nil {
		log.Printf("error running travellers import: %v", err)
		return
	}
}

func profileRun() {
	var buf bytes.Buffer

	if err := pprof.StartCPUProfile(&buf); err != nil {
		log.Fatal(err)
		return
	}

	log.Println("CPU profile started")

	concurrentRun()

	pprof.StopCPUProfile()

	log.Println("CPU profile ended")

	f, err := os.Create("cpu_profile.prof")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()

	if _, err = f.Write(buf.Bytes()); err != nil {
		log.Fatal(err)
		return
	}
}

func traceRun() {
	var buf bytes.Buffer

	if err := trace.Start(&buf); err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Trace started")

	//simpleRun()
	concurrentRun()

	trace.Stop()

	log.Println("Trace ended")

	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()

	if _, err = f.Write(buf.Bytes()); err != nil {
		log.Fatal(err)
		return
	}
}

//func traceRun() {
//	fr := trace.NewFlightRecorder()
//	if err := fr.Start(); err != nil {
//		// handle error
//	}
//
//	defer func() {
//		if err := fr.Stop(); err != nil {
//			// handle error
//		}
//	}()
//
//	app := fiber.New()
//
//	app.Get("/trace", func(ctx fiber.Ctx) error {
//		f, err := os.OpenFile("file.out", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
//		if err != nil {
//			// handle error
//		}
//
//		if _, err := fr.WriteTo(f); err != nil {
//			// handle error
//		}
//		return ctx.SendStatus(fiber.StatusOK)
//	})
//
//	dbExec, err := newDB()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	ctx := context.Background()
//
//	travellersClient := postgres.NewClient(dbExec)
//	travellersService := internal.NewTravellers(travellersClient)
//	travellersParser := ftp.NewParser(travellersService)
//
//	if err = travellersParser.Run(ctx, "/Users/denys/Go/src/github.com/DenisGoldiner/webapp/internal/integration/data/test_1.csv"); err != nil {
//		log.Printf("error running travellers import: %v", err)
//		return
//	}
//
//	if err := app.Listen(":8080"); err != nil {
//		// handle error
//		app.Shutdown()
//	}
//}

func newDB() (sqlx.ExtContext, error) {
	dsn := "postgres://postgres:postgres@localhost:5432/travellers?sslmode=disable"
	conn, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
