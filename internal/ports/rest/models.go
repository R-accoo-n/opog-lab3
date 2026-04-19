package rest

import (
	"strings"

	"github.com/DenisGoldiner/webapp/internal"
	"github.com/google/uuid"
)

type Traveller struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Age       int       `json:"-"`
}

// CreateTravellerPayload represents the payload for creating a new traveller.
// Example:
//
//	{
//	  "first_name": "John",
//	  "last_name": "Doe",
//	  "age": 30
//	}
type CreateTravellerPayload struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}

func (p CreateTravellerPayload) toServiceParams() internal.CreateTravellerPayload {
	return internal.CreateTravellerPayload{
		FirstName: strings.TrimSpace(p.FirstName),
		LastName:  strings.TrimSpace(p.LastName),
		Age:       p.Age,
	}
}

type CreateTravellerResponse struct {
	ID uuid.UUID `json:"id"`
}
