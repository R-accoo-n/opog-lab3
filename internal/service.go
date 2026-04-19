package internal

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Travellers struct {
	travellerStorage TravellerStorage
}

func NewTravellers(db TravellerStorage) Travellers {
	return Travellers{travellerStorage: db}
}

func (t Travellers) GetTraveller(ctx context.Context, id uuid.UUID) (Traveller, error) {
	if id == uuid.Nil {
		return Traveller{}, fmt.Errorf("id must be a valid uuid")
	}

	res, err := t.travellerStorage.Get(ctx, id)
	if err != nil {
		return Traveller{}, fmt.Errorf("%w: failed to get traveller from travellerStorage", err)
	}

	return res, nil
}

func (t Travellers) CreateTraveller(ctx context.Context, traveller CreateTravellerPayload) (uuid.UUID, error) {
	if traveller.FirstName == "" || traveller.LastName == "" {
		return uuid.Nil, fmt.Errorf("%w: first name and last name must be provided", ErrInvalidInput)
	}

	travellerID, err := t.travellerStorage.Create(ctx, traveller)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create traveller in travellerStorage: %w", err)
	}

	return travellerID, nil
}

func (t Travellers) BulkCreateTravellers(ctx context.Context, travellers []CreateTravellerPayload) ([]uuid.UUID, error) {
	if len(travellers) == 0 {
		return nil, fmt.Errorf("%w: travellers list must not be empty", ErrInvalidInput)
	}

	travellerIDs, err := t.travellerStorage.BulkCreate(ctx, travellers)
	if err != nil {
		return nil, fmt.Errorf("failed to bulk create travellers in travellerStorage: %w", err)
	}

	return travellerIDs, nil
}

func (t Travellers) DeleteTraveller() {

}
