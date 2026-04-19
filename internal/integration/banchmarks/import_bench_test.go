package banchmarks

import (
	"context"
	"log"
	"testing"

	"github.com/DenisGoldiner/webapp/internal"
	"github.com/DenisGoldiner/webapp/internal/adapters/postgres"
	"github.com/DenisGoldiner/webapp/internal/ports/ftp"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// simple	  BenchmarkSample-12    	       4	 274281636 ns/op
// concurrent BenchmarkSample-12    	       4	 269966636 ns/op
//			  BenchmarkSample-12    	       4	 264737292 ns/op
//            BenchmarkSample-12    	       5	 245282900 ns/op

// BenchmarkSample-12    	       1	2764022250 ns/op
// BenchmarkSample-12    	       2	 747799854 ns/op
// concurrent 10 * 50 					BenchmarkSample-12    	      25	  47567932 ns/op
// simple all in 						BenchmarkSample-12    	      18	  60528782 ns/op
// concurrent 5 * 500 					BenchmarkSample-12    	      39	  26125076 ns/op
func BenchmarkSample(b *testing.B) {
	dbExec, err := newDB()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	travellersClient := postgres.NewClient(dbExec)
	travellersService := internal.NewTravellers(travellersClient)
	travellersParser := ftp.NewConcurrentParser(travellersService)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err = travellersParser.Run(ctx, "/Users/denys/Go/src/github.com/DenisGoldiner/webapp/internal/integration/data/test_1.csv"); err != nil {
			b.Fatalf("unespected error, %v", err)
		}
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
