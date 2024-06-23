// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package schema

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID
	FirstName       string
	LastName        string
	Email           string
	Password        []byte
	EmailVerifiedAt sql.NullTime
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
