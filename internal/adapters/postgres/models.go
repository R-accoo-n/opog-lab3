package postgres

import "github.com/google/uuid"

const pqErrCodeUniqueViolation = "23505"

type Traveller struct {
	ID        uuid.UUID `db:"id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Age       int       `db:"age"`
}
