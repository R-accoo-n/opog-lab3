package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/DenisGoldiner/webapp/internal"
)

type Client struct {
	dbExec sqlx.ExtContext
}

func NewClient(dbExec sqlx.ExtContext) Client {
	return Client{dbExec: dbExec}
}

func (c Client) Get(ctx context.Context, id uuid.UUID) (internal.Traveller, error) {
	q := "select id, first_name, last_name, age from travellers where id = $1"

	rows, err := c.dbExec.QueryxContext(ctx, q, id)
	if err != nil {
		return internal.Traveller{}, fmt.Errorf("failed to fetch traveler: %w", err)
	}

	defer func() { _ = rows.Close() }()

	var travelers []Traveller

	for rows.Next() {
		var traveler Traveller

		if err = rows.StructScan(&traveler); err != nil {
			return internal.Traveller{}, fmt.Errorf("failed to scan traveler: %w", err)
		}

		travelers = append(travelers, traveler)
	}

	if len(travelers) == 0 {
		return internal.Traveller{}, fmt.Errorf("no travelers with id %s: %w", id, internal.ErrNoResource)
	}

	return internal.Traveller{
		ID:        travelers[0].ID,
		FirstName: travelers[0].FirstName,
		LastName:  travelers[0].LastName,
		Age:       travelers[0].Age,
	}, err
}

func (c Client) Select(ctx context.Context, ids []uuid.UUID, limit, offset int) ([]internal.Traveller, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	q := "select id, first_name, last_name, age from travellers "
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))

	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	q += fmt.Sprintf(" where id IN (%s) ", strings.Join(placeholders, ", "))
	q += fmt.Sprintf(" limit %d offset %d ", limit, offset)
	q += " order by created_at desc "

	rows, err := c.dbExec.QueryxContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch traveler: %w", err)
	}

	defer func() { _ = rows.Close() }()

	var travelers []Traveller

	for rows.Next() {
		var traveler Traveller

		if err = rows.StructScan(&traveler); err != nil {
			return nil, fmt.Errorf("failed to scan traveler: %w", err)
		}

		travelers = append(travelers, traveler)
	}

	if len(travelers) == 0 {
		return nil, fmt.Errorf("no travelers", internal.ErrNoResource)
	}

	var internalTravelers []internal.Traveller
	for i := range travelers {
		internalTravelers = append(internalTravelers, internal.Traveller{
			ID:        travelers[i].ID,
			FirstName: travelers[i].FirstName,
			LastName:  travelers[i].LastName,
			Age:       travelers[i].Age,
		})
	}

	return internalTravelers, err
}

func (c Client) Create(ctx context.Context, params internal.CreateTravellerPayload) (uuid.UUID, error) {
	q := "insert into travellers (first_name, last_name, age) values ($1, $2, $3) returning id"

	rows, err := c.dbExec.QueryxContext(ctx, q, params.FirstName, params.LastName, params.Age)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pqErrCodeUniqueViolation {
			return uuid.Nil, fmt.Errorf("traveler already exists: %w", internal.ErrAlreadyExists)
		}

		return uuid.Nil, fmt.Errorf("failed to create traveler: %w", err)
	}

	defer func() { _ = rows.Close() }()

	var travelerIDs []uuid.UUID

	for rows.Next() {
		var travelerID uuid.UUID

		if err = rows.Scan(&travelerID); err != nil {
			return uuid.Nil, fmt.Errorf("failed to scan traveler id: %w", err)
		}

		travelerIDs = append(travelerIDs, travelerID)
	}

	if len(travelerIDs) == 0 {
		return uuid.Nil, fmt.Errorf("failed to create traveler: no id returned")
	}

	return travelerIDs[0], nil
}

func (c Client) BulkCreate(ctx context.Context, params []internal.CreateTravellerPayload) ([]uuid.UUID, error) {
	if len(params) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(params))
	args := make([]interface{}, 0, len(params)*3)

	for i, p := range params {
		base := i * 3
		placeholders[i] = fmt.Sprintf("($%d, $%d, $%d)", base+1, base+2, base+3)
		args = append(args, p.FirstName, p.LastName, p.Age)
	}

	q := "INSERT INTO travellers (first_name, last_name, age) VALUES " +
		strings.Join(placeholders, ", ") +
		" RETURNING id"

	rows, err := c.dbExec.QueryxContext(ctx, q, args...)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pqErrCodeUniqueViolation {
			return nil, fmt.Errorf("traveler already exists: %w", internal.ErrAlreadyExists)
		}

		return nil, fmt.Errorf("failed to bulk create travelers: %w", err)
	}

	defer func() { _ = rows.Close() }()

	var ids []uuid.UUID

	for rows.Next() {
		var id uuid.UUID

		if err = rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan traveler id: %w", err)
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (c Client) DeleteTraveller() {

}
