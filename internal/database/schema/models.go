// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package schema

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Exam struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	Title            string
	Slug             pgtype.Text
	VisibilityStatus string
	IsAccessable     pgtype.Bool
	CreatedAt        pgtype.Timestamp
	UpdatedAt        pgtype.Timestamp
	LiveStatus       pgtype.Text
}

type ExamParticipant struct {
	UserID           pgtype.UUID
	ExamID           uuid.UUID
	Email            string
	Status           string
	CreatedAt        pgtype.Timestamp
	UpdatedAt        pgtype.Timestamp
	AttendanceStatus pgtype.Text
}

type ExamQuestion struct {
	ID        uuid.UUID
	ExamID    uuid.UUID
	Question  string
	Answers   []byte
	Type      string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type QuestionAnswer struct {
	UserID     uuid.UUID
	QuestionID uuid.UUID
	Answer     string
	IsCorrect  pgtype.Bool
}

type Token struct {
	Hash   []byte
	UserID uuid.UUID
	Expiry pgtype.Timestamp
	Scope  string
}

type User struct {
	ID              uuid.UUID
	FirstName       string
	LastName        string
	Email           string
	Password        []byte
	EmailVerifiedAt pgtype.Timestamp
	CreatedAt       pgtype.Timestamp
	UpdatedAt       pgtype.Timestamp
}
