package internal

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrNoResource = errors.New("no resource found")
var ErrAlreadyExists = errors.New("resource already exists")
var ErrInvalidInput = errors.New("invalid input")

type TravellerStorage interface {
	Get(ctx context.Context, id uuid.UUID) (Traveller, error)
	Create(ctx context.Context, params CreateTravellerPayload) (uuid.UUID, error)
	BulkCreate(ctx context.Context, params []CreateTravellerPayload) ([]uuid.UUID, error)
}

type Traveller struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Age       int
}

type CreateTravellerPayload struct {
	FirstName string
	LastName  string
	Age       int
}
