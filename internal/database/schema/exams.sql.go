// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: exams.sql

package schema

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
)

const checkExamTitleExists = `-- name: CheckExamTitleExists :one
SELECT EXISTS(SELECT 1 FROM exams WHERE title = $1 AND user_id = $2)
`

type CheckExamTitleExistsParams struct {
	Title  string
	UserID uuid.UUID
}

func (q *Queries) CheckExamTitleExists(ctx context.Context, arg CheckExamTitleExistsParams) (bool, error) {
	row := q.db.QueryRowContext(ctx, checkExamTitleExists, arg.Title, arg.UserID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createExam = `-- name: CreateExam :one
INSERT INTO exams(id , user_id , title , slug , visibility_status , is_accessable,  created_at , updated_at)
VALUES (uuid_generate_v4() , $1 , $2 , $3 , $4 ,$5 ,  current_timestamp , current_timestamp)
RETURNING id, user_id, title, slug, visibility_status, is_accessable, created_at, updated_at
`

type CreateExamParams struct {
	UserID           uuid.UUID
	Title            string
	Slug             sql.NullString
	VisibilityStatus string
	IsAccessable     sql.NullBool
}

func (q *Queries) CreateExam(ctx context.Context, arg CreateExamParams) (Exam, error) {
	row := q.db.QueryRowContext(ctx, createExam,
		arg.UserID,
		arg.Title,
		arg.Slug,
		arg.VisibilityStatus,
		arg.IsAccessable,
	)
	var i Exam
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.Slug,
		&i.VisibilityStatus,
		&i.IsAccessable,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createExamQuestions = `-- name: CreateExamQuestions :one
INSERT INTO exam_questions(id ,  exam_id , question , type , answers ,  created_at , updated_at)
VALUES (uuid_generate_v4() , $1 , $2 , $3 , $4 , current_timestamp , current_timestamp)
RETURNING id, exam_id, question, answers, type, created_at, updated_at
`

type CreateExamQuestionsParams struct {
	ExamID   uuid.UUID
	Question string
	Type     string
	Answers  json.RawMessage
}

func (q *Queries) CreateExamQuestions(ctx context.Context, arg CreateExamQuestionsParams) (ExamQuestion, error) {
	row := q.db.QueryRowContext(ctx, createExamQuestions,
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

const getExamById = `-- name: GetExamById :one
SELECT id, user_id, title, slug, visibility_status, is_accessable, created_at, updated_at FROM exams WHERE id = $1
`

func (q *Queries) GetExamById(ctx context.Context, id uuid.UUID) (Exam, error) {
	row := q.db.QueryRowContext(ctx, getExamById, id)
	var i Exam
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.Slug,
		&i.VisibilityStatus,
		&i.IsAccessable,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserExams = `-- name: GetUserExams :many
SELECT id, user_id, title, slug, visibility_status, is_accessable, created_at, updated_at FROM exams WHERE user_id = $1
`

func (q *Queries) GetUserExams(ctx context.Context, userID uuid.UUID) ([]Exam, error) {
	rows, err := q.db.QueryContext(ctx, getUserExams, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Exam
	for rows.Next() {
		var i Exam
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.Slug,
			&i.VisibilityStatus,
			&i.IsAccessable,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
