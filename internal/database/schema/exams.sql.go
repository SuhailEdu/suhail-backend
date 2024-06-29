// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: exams.sql

package schema

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const checkExamTitleExists = `-- name: CheckExamTitleExists :one
SELECT EXISTS(SELECT 1 FROM exams WHERE title = $1 AND user_id = $2)
`

type CheckExamTitleExistsParams struct {
	Title  string
	UserID uuid.UUID
}

func (q *Queries) CheckExamTitleExists(ctx context.Context, arg CheckExamTitleExistsParams) (bool, error) {
	row := q.db.QueryRow(ctx, checkExamTitleExists, arg.Title, arg.UserID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createExam = `-- name: CreateExam :one
INSERT INTO exams(id, user_id, title, slug, visibility_status, is_accessable, created_at, updated_at)
VALUES (uuid_generate_v4(), $1, $2, $3, $4, $5, current_timestamp, current_timestamp)
RETURNING id, user_id, title, slug, visibility_status, is_accessable, created_at, updated_at
`

type CreateExamParams struct {
	UserID           uuid.UUID
	Title            string
	Slug             pgtype.Text
	VisibilityStatus string
	IsAccessable     pgtype.Bool
}

func (q *Queries) CreateExam(ctx context.Context, arg CreateExamParams) (Exam, error) {
	row := q.db.QueryRow(ctx, createExam,
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

const deleteExam = `-- name: DeleteExam :exec
DELETE
FROM exam_questions
WHERE id = $1
`

func (q *Queries) DeleteExam(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteExam, id)
	return err
}

const findMyExam = `-- name: FindMyExam :many
SELECT exams.id, exams.user_id, exams.title, exams.slug, exams.visibility_status, exams.is_accessable, exams.created_at, exams.updated_at, exam_questions.id, exam_questions.exam_id, exam_questions.question, exam_questions.answers, exam_questions.type, exam_questions.created_at, exam_questions.updated_at
FROM exams
         LEFT JOIN exam_questions ON exam_questions.exam_id = exams.id
WHERE exams.id = $1
  AND exams.user_id = $2
`

type FindMyExamParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

type FindMyExamRow struct {
	Exam         Exam
	ExamQuestion ExamQuestion
}

func (q *Queries) FindMyExam(ctx context.Context, arg FindMyExamParams) ([]FindMyExamRow, error) {
	rows, err := q.db.Query(ctx, findMyExam, arg.ID, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindMyExamRow
	for rows.Next() {
		var i FindMyExamRow
		if err := rows.Scan(
			&i.Exam.ID,
			&i.Exam.UserID,
			&i.Exam.Title,
			&i.Exam.Slug,
			&i.Exam.VisibilityStatus,
			&i.Exam.IsAccessable,
			&i.Exam.CreatedAt,
			&i.Exam.UpdatedAt,
			&i.ExamQuestion.ID,
			&i.ExamQuestion.ExamID,
			&i.ExamQuestion.Question,
			&i.ExamQuestion.Answers,
			&i.ExamQuestion.Type,
			&i.ExamQuestion.CreatedAt,
			&i.ExamQuestion.UpdatedAt,
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

const findMyParticipatedExam = `-- name: FindMyParticipatedExam :many
SELECT exams.id, exams.user_id, exams.title, exams.slug, exams.visibility_status, exams.is_accessable, exams.created_at, exams.updated_at, exam_questions.id, exam_questions.exam_id, exam_questions.question, exam_questions.answers, exam_questions.type, exam_questions.created_at, exam_questions.updated_at
FROM exams
         INNER JOIN exam_participants ON exam_participants.exam_id = exams.id
         LEFT JOIN exam_questions ON exam_questions.exam_id = exams.id
WHERE exam_participants.user_id = $1
  AND exams.id = $2
`

type FindMyParticipatedExamParams struct {
	UserID uuid.UUID
	ID     uuid.UUID
}

type FindMyParticipatedExamRow struct {
	Exam         Exam
	ExamQuestion ExamQuestion
}

func (q *Queries) FindMyParticipatedExam(ctx context.Context, arg FindMyParticipatedExamParams) ([]FindMyParticipatedExamRow, error) {
	rows, err := q.db.Query(ctx, findMyParticipatedExam, arg.UserID, arg.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindMyParticipatedExamRow
	for rows.Next() {
		var i FindMyParticipatedExamRow
		if err := rows.Scan(
			&i.Exam.ID,
			&i.Exam.UserID,
			&i.Exam.Title,
			&i.Exam.Slug,
			&i.Exam.VisibilityStatus,
			&i.Exam.IsAccessable,
			&i.Exam.CreatedAt,
			&i.Exam.UpdatedAt,
			&i.ExamQuestion.ID,
			&i.ExamQuestion.ExamID,
			&i.ExamQuestion.Question,
			&i.ExamQuestion.Answers,
			&i.ExamQuestion.Type,
			&i.ExamQuestion.CreatedAt,
			&i.ExamQuestion.UpdatedAt,
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

const getExamById = `-- name: GetExamById :one
SELECT id, user_id, title, slug, visibility_status, is_accessable, created_at, updated_at
FROM exams
WHERE id = $1
`

func (q *Queries) GetExamById(ctx context.Context, id uuid.UUID) (Exam, error) {
	row := q.db.QueryRow(ctx, getExamById, id)
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

const getExamQuestions = `-- name: GetExamQuestions :many
SELECT id, exam_id, question, answers, type, created_at, updated_at
FROM exam_questions
WHERE exam_id = $1
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

const getParticipatedExams = `-- name: GetParticipatedExams :many
SELECT exams.id, exams.user_id, exams.title, exams.slug, exams.visibility_status, exams.is_accessable, exams.created_at, exams.updated_at, COUNT(exam_questions.*) as questions_count
FROM exams
         LEFT JOIN exam_questions ON exam_questions.exam_id = exams.id
         INNER JOIN exam_participants ON exam_participants.exam_id = exams.id AND exam_participants.user_id = $1
GROUP BY exams.id
ORDER BY exams.created_at DESC
`

type GetParticipatedExamsRow struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	Title            string
	Slug             pgtype.Text
	VisibilityStatus string
	IsAccessable     pgtype.Bool
	CreatedAt        pgtype.Timestamp
	UpdatedAt        pgtype.Timestamp
	QuestionsCount   int64
}

// WHERE user_id = $1
func (q *Queries) GetParticipatedExams(ctx context.Context, userID uuid.UUID) ([]GetParticipatedExamsRow, error) {
	rows, err := q.db.Query(ctx, getParticipatedExams, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetParticipatedExamsRow
	for rows.Next() {
		var i GetParticipatedExamsRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.Slug,
			&i.VisibilityStatus,
			&i.IsAccessable,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.QuestionsCount,
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

const getUserExams = `-- name: GetUserExams :many
SELECT exams.id, exams.user_id, exams.title, exams.slug, exams.visibility_status, exams.is_accessable, exams.created_at, exams.updated_at, COUNT(exam_participants.*) as particpants_count, COUNT(exam_questions.*) as questions_count
FROM exams
         LEFT JOIN exam_questions ON exam_questions.exam_id = exams.id
         LEFT JOIN exam_participants ON exam_participants.exam_id = exams.id
WHERE exams.user_id = $1
GROUP BY exams.id
ORDER BY exams.created_at DESC
`

type GetUserExamsRow struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	Title            string
	Slug             pgtype.Text
	VisibilityStatus string
	IsAccessable     pgtype.Bool
	CreatedAt        pgtype.Timestamp
	UpdatedAt        pgtype.Timestamp
	ParticpantsCount int64
	QuestionsCount   int64
}

func (q *Queries) GetUserExams(ctx context.Context, userID uuid.UUID) ([]GetUserExamsRow, error) {
	rows, err := q.db.Query(ctx, getUserExams, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserExamsRow
	for rows.Next() {
		var i GetUserExamsRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.Slug,
			&i.VisibilityStatus,
			&i.IsAccessable,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ParticpantsCount,
			&i.QuestionsCount,
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

const getUserExamsWithQuestions = `-- name: GetUserExamsWithQuestions :many
SELECT exams.id, user_id, title, slug, visibility_status, is_accessable, exams.created_at, exams.updated_at, exam_questions.id, exam_id, question, answers, type, exam_questions.created_at, exam_questions.updated_at
FROM exams
         LEFT JOIN exam_questions ON exam_questions.exam_id = exams.id
WHERE user_id = $1
`

type GetUserExamsWithQuestionsRow struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	Title            string
	Slug             pgtype.Text
	VisibilityStatus string
	IsAccessable     pgtype.Bool
	CreatedAt        pgtype.Timestamp
	UpdatedAt        pgtype.Timestamp
	ID_2             pgtype.UUID
	ExamID           pgtype.UUID
	Question         pgtype.Text
	Answers          []byte
	Type             pgtype.Text
	CreatedAt_2      pgtype.Timestamp
	UpdatedAt_2      pgtype.Timestamp
}

func (q *Queries) GetUserExamsWithQuestions(ctx context.Context, userID uuid.UUID) ([]GetUserExamsWithQuestionsRow, error) {
	rows, err := q.db.Query(ctx, getUserExamsWithQuestions, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserExamsWithQuestionsRow
	for rows.Next() {
		var i GetUserExamsWithQuestionsRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.Slug,
			&i.VisibilityStatus,
			&i.IsAccessable,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ID_2,
			&i.ExamID,
			&i.Question,
			&i.Answers,
			&i.Type,
			&i.CreatedAt_2,
			&i.UpdatedAt_2,
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

const updateExam = `-- name: UpdateExam :exec
UPDATE exams
SET title             = $1,
    visibility_status = $2
WHERE id = $3

RETURNING id, user_id, title, slug, visibility_status, is_accessable, created_at, updated_at
`

type UpdateExamParams struct {
	Title            string
	VisibilityStatus string
	ID               uuid.UUID
}

func (q *Queries) UpdateExam(ctx context.Context, arg UpdateExamParams) error {
	_, err := q.db.Exec(ctx, updateExam, arg.Title, arg.VisibilityStatus, arg.ID)
	return err
}
