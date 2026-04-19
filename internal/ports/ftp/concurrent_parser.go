package ftp

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/DenisGoldiner/webapp/internal"
)

type ConcurrentParser struct {
	service internal.Travellers
}

func NewConcurrentParser(service internal.Travellers) ConcurrentParser {
	return ConcurrentParser{service: service}
}

func (p ConcurrentParser) Run(ctx context.Context, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open the file %s: %w", filePath, err)
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	reader := csv.NewReader(bytes.NewReader(data))

	const workers = 5
	ch := p.parse(reader, workers)

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Go(func() {
			//p.process(ctx, ch)
			p.bulkProcess(ctx, ch)
		})
	}

	wg.Wait()

	return nil
}

func (p ConcurrentParser) parse(r *csv.Reader, workers int) <-chan internal.CreateTravellerPayload {
	ch := make(chan internal.CreateTravellerPayload, 10*workers)

	go func() {
		for i := 0; ; i++ {
			row, err := r.Read()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Printf("failed to parse row %d: %w", i, err)
				continue
			}

			age, err := strconv.Atoi(row[2])
			if err != nil {
				log.Printf("failed to parse age %d: %w", i, err)
				continue
			}

			traveler := internal.CreateTravellerPayload{
				FirstName: strings.TrimSpace(row[0]),
				LastName:  strings.TrimSpace(row[1]),
				Age:       age,
			}

			ch <- traveler
		}

		close(ch)
	}()

	return ch
}

func (p ConcurrentParser) process(ctx context.Context, ch <-chan internal.CreateTravellerPayload) {
	for traveler := range ch {
		if _, err := p.service.CreateTraveller(ctx, traveler); err != nil {
			log.Printf("failed to create traveller: %v", err)
			continue
		}
	}
}

func (p ConcurrentParser) bulkProcess(ctx context.Context, ch <-chan internal.CreateTravellerPayload) {
	bulk := make([]internal.CreateTravellerPayload, 0, 50)

	for traveler := range ch {
		if len(bulk) < 500 {
			bulk = append(bulk, traveler)
			continue
		}

		if _, err := p.service.BulkCreateTravellers(ctx, bulk); err != nil {
			log.Printf("failed to create traveller: %v", err)
			continue
		}

		bulk = bulk[:0]
	}
}
