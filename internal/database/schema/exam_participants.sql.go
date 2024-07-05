// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: exam_participants.sql

package schema

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const checkParticipant = `-- name: CheckParticipant :one
SELECT EXISTS(SELECT 1 FROM exam_participants WHERE exam_id = $1 AND user_id = $2)
`

type CheckParticipantParams struct {
	ExamID uuid.UUID
	UserID pgtype.UUID
}

func (q *Queries) CheckParticipant(ctx context.Context, arg CheckParticipantParams) (bool, error) {
	row := q.db.QueryRow(ctx, checkParticipant, arg.ExamID, arg.UserID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkParticipantByToken = `-- name: CheckParticipantByToken :one
SELECT EXISTS(SELECT 1
              FROM exam_participants
              WHERE exam_id = $1
                AND user_id = (SELECT user_id
                               FROM tokens
                               WHERE tokens.hash = $2
                               LIMIT 1))
`

type CheckParticipantByTokenParams struct {
	ExamID uuid.UUID
	Hash   []byte
}

func (q *Queries) CheckParticipantByToken(ctx context.Context, arg CheckParticipantByTokenParams) (bool, error) {
	row := q.db.QueryRow(ctx, checkParticipantByToken, arg.ExamID, arg.Hash)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createExamParticipant = `-- name: CreateExamParticipant :exec
INSERT INTO exam_participants (exam_id, email, status)
VALUES ($1, $2, $3)
ON CONFLICT (email , exam_id) DO NOTHING
`

type CreateExamParticipantParams struct {
	ExamID uuid.UUID
	Email  string
	Status string
}

func (q *Queries) CreateExamParticipant(ctx context.Context, arg CreateExamParticipantParams) error {
	_, err := q.db.Exec(ctx, createExamParticipant, arg.ExamID, arg.Email, arg.Status)
	return err
}

const deleteParticipants = `-- name: DeleteParticipants :exec
DELETE
FROM exam_participants
WHERE email = ANY ($2)
  AND exam_id = $1
`

type DeleteParticipantsParams struct {
	ExamID uuid.UUID
	Emails []string
}

func (q *Queries) DeleteParticipants(ctx context.Context, arg DeleteParticipantsParams) error {
	_, err := q.db.Exec(ctx, deleteParticipants, arg.ExamID, arg.Emails)
	return err
}

const getExamParticipants = `-- name: GetExamParticipants :many
SELECT exam_participants.email, exam_participants.status, users.id, users.first_name, users.last_name, users.email, users.password, users.email_verified_at, users.created_at, users.updated_at
FROM exam_participants
         LEFT JOIN users on users.id = exam_participants.user_id
WHERE exam_id = $1
`

type GetExamParticipantsRow struct {
	Email           string
	Status          string
	ID              pgtype.UUID
	FirstName       pgtype.Text
	LastName        pgtype.Text
	Email_2         pgtype.Text
	Password        []byte
	EmailVerifiedAt pgtype.Timestamp
	CreatedAt       pgtype.Timestamp
	UpdatedAt       pgtype.Timestamp
}

func (q *Queries) GetExamParticipants(ctx context.Context, examID uuid.UUID) ([]GetExamParticipantsRow, error) {
	rows, err := q.db.Query(ctx, getExamParticipants, examID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetExamParticipantsRow
	for rows.Next() {
		var i GetExamParticipantsRow
		if err := rows.Scan(
			&i.Email,
			&i.Status,
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email_2,
			&i.Password,
			&i.EmailVerifiedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
