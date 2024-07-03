// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: exam_questions.sql

package schema

import (
	"context"

	"github.com/google/uuid"
)

const checkQuestionExits = `-- name: CheckQuestionExits :one
SELECT EXISTS(SELECT 1 FROM exam_questions WHERE id = $1 AND exam_id = $2)
`

type CheckQuestionExitsParams struct {
	ID     uuid.UUID
	ExamID uuid.UUID
}

func (q *Queries) CheckQuestionExits(ctx context.Context, arg CheckQuestionExitsParams) (bool, error) {
	row := q.db.QueryRow(ctx, checkQuestionExits, arg.ID, arg.ExamID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkQuestionTitleExits = `-- name: CheckQuestionTitleExits :one
SELECT EXISTS(SELECT 1 FROM exam_questions WHERE question = $1 AND exam_id = $2)
`

type CheckQuestionTitleExitsParams struct {
	Question string
	ExamID   uuid.UUID
}

func (q *Queries) CheckQuestionTitleExits(ctx context.Context, arg CheckQuestionTitleExitsParams) (bool, error) {
	row := q.db.QueryRow(ctx, checkQuestionTitleExits, arg.Question, arg.ExamID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

type CreateExamQuestionsParams struct {
	ExamID   uuid.UUID
	Question string
	Type     string
	Answers  []byte
}

const createQuestion = `-- name: CreateQuestion :one
INSERT INTO exam_questions(exam_id, question, type, answers)
VALUES ($1, $2, $3, $4)
RETURNING id, exam_id, question, answers, type, created_at, updated_at
`

type CreateQuestionParams struct {
	ExamID   uuid.UUID
	Question string
	Type     string
	Answers  []byte
}

func (q *Queries) CreateQuestion(ctx context.Context, arg CreateQuestionParams) (ExamQuestion, error) {
	row := q.db.QueryRow(ctx, createQuestion,
		arg.ExamID,
		arg.Question,
		arg.Type,
		arg.Answers,
	)
	var i ExamQuestion
	err := row.Scan(
		&i.ID,
		&i.ExamID,
		&i.Question,
		&i.Answers,
		&i.Type,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteQuestion = `-- name: DeleteQuestion :exec
DELETE
FROM exam_questions

WHERE exam_questions.id = $1
  AND EXISTS(SELECT 1 FROM exams WHERE exams.user_id = $2 AND exams.id = $3)
`

type DeleteQuestionParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
	ID_2   uuid.UUID
}

func (q *Queries) DeleteQuestion(ctx context.Context, arg DeleteQuestionParams) error {
	_, err := q.db.Exec(ctx, deleteQuestion, arg.ID, arg.UserID, arg.ID_2)
	return err
}

const getExamQuestions = `-- name: GetExamQuestions :many
SELECT id, exam_id, question, answers, type, created_at, updated_at
FROM exam_questions
WHERE exam_id = $1
LIMIT 1
`

func (q *Queries) GetExamQuestions(ctx context.Context, examID uuid.UUID) ([]ExamQuestion, error) {
	rows, err := q.db.Query(ctx, getExamQuestions, examID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ExamQuestion
	for rows.Next() {
		var i ExamQuestion
		if err := rows.Scan(
			&i.ID,
			&i.ExamID,
			&i.Question,
			&i.Answers,
			&i.Type,
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

const getQuestionById = `-- name: GetQuestionById :one
SELECT id, exam_id, question, answers, type, created_at, updated_at
FROM exam_questions
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetQuestionById(ctx context.Context, id uuid.UUID) (ExamQuestion, error) {
	row := q.db.QueryRow(ctx, getQuestionById, id)
	var i ExamQuestion
	err := row.Scan(
		&i.ID,
		&i.ExamID,
		&i.Question,
		&i.Answers,
		&i.Type,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateQuestion = `-- name: UpdateQuestion :exec
UPDATE exam_questions
SET question = $1,
    answers  = $2
WHERE id = $3

RETURNING id, exam_id, question, answers, type, created_at, updated_at
`

type UpdateQuestionParams struct {
	Question string
	Answers  []byte
	ID       uuid.UUID
}

func (q *Queries) UpdateQuestion(ctx context.Context, arg UpdateQuestionParams) error {
	_, err := q.db.Exec(ctx, updateQuestion, arg.Question, arg.Answers, arg.ID)
	return err
}
