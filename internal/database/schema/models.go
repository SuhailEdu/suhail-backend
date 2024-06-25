// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package schema

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Exam struct {
	ID               pgtype.UUID
	UserID           pgtype.UUID
	Title            string
	Slug             pgtype.Text
	VisibilityStatus string
	IsAccessable     pgtype.Bool
	CreatedAt        pgtype.Timestamp
	UpdatedAt        pgtype.Timestamp
}

type ExamParticipant struct {
	UserID    pgtype.UUID
	ExamID    pgtype.UUID
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type ExamQuestion struct {
	ID        pgtype.UUID
	ExamID    pgtype.UUID
	Question  string
	Answers   []byte
	Type      string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type Token struct {
	Hash   []byte
	UserID pgtype.UUID
	Expiry pgtype.Timestamp
	Scope  string
}

type User struct {
	ID              pgtype.UUID
	FirstName       string
	LastName        string
	Email           string
	Password        []byte
	EmailVerifiedAt pgtype.Timestamp
	CreatedAt       pgtype.Timestamp
	UpdatedAt       pgtype.Timestamp
}
