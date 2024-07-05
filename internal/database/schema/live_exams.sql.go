// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: live_exams.sql

package schema

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const getLiveExamParticipants = `-- name: GetLiveExamParticipants :many
SELECT users.id, users.first_name, users.last_name, users.email, users.password, users.email_verified_at, users.created_at, users.updated_at, exam_participants.user_id, exam_participants.exam_id, exam_participants.email, exam_participants.status, exam_participants.created_at, exam_participants.updated_at, exam_participants.attendance_status
FROM exam_participants
         INNER JOIN users ON users.id = exam_participants.user_id
WHERE exam_participants.exam_id = $1
`

type GetLiveExamParticipantsRow struct {
	User            User
	ExamParticipant ExamParticipant
}

func (q *Queries) GetLiveExamParticipants(ctx context.Context, examID uuid.UUID) ([]GetLiveExamParticipantsRow, error) {
	rows, err := q.db.Query(ctx, getLiveExamParticipants, examID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetLiveExamParticipantsRow
	for rows.Next() {
		var i GetLiveExamParticipantsRow
		if err := rows.Scan(
			&i.User.ID,
			&i.User.FirstName,
			&i.User.LastName,
			&i.User.Email,
			&i.User.Password,
			&i.User.EmailVerifiedAt,
			&i.User.CreatedAt,
			&i.User.UpdatedAt,
			&i.ExamParticipant.UserID,
			&i.ExamParticipant.ExamID,
			&i.ExamParticipant.Email,
			&i.ExamParticipant.Status,
			&i.ExamParticipant.CreatedAt,
			&i.ExamParticipant.UpdatedAt,
			&i.ExamParticipant.AttendanceStatus,
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

const getLiveExamQuestionForManager = `-- name: GetLiveExamQuestionForManager :many
SELECT exam_questions.id, exam_questions.exam_id, exam_questions.question, exam_questions.answers, exam_questions.type, exam_questions.created_at, exam_questions.updated_at, exams.id, exams.user_id, exams.title, exams.slug, exams.visibility_status, exams.is_accessable, exams.created_at, exams.updated_at, exams.live_status, exams.ip_range_start, exams.ip_range_end
FROM exam_questions
         INNER JOIN exams ON exams.id = $1
LIMIT 1
`

type GetLiveExamQuestionForManagerRow struct {
	ExamQuestion ExamQuestion
	Exam         Exam
}

func (q *Queries) GetLiveExamQuestionForManager(ctx context.Context, id uuid.UUID) ([]GetLiveExamQuestionForManagerRow, error) {
	rows, err := q.db.Query(ctx, getLiveExamQuestionForManager, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetLiveExamQuestionForManagerRow
	for rows.Next() {
		var i GetLiveExamQuestionForManagerRow
		if err := rows.Scan(
			&i.ExamQuestion.ID,
			&i.ExamQuestion.ExamID,
			&i.ExamQuestion.Question,
			&i.ExamQuestion.Answers,
			&i.ExamQuestion.Type,
			&i.ExamQuestion.CreatedAt,
			&i.ExamQuestion.UpdatedAt,
			&i.Exam.ID,
			&i.Exam.UserID,
			&i.Exam.Title,
			&i.Exam.Slug,
			&i.Exam.VisibilityStatus,
			&i.Exam.IsAccessable,
			&i.Exam.CreatedAt,
			&i.Exam.UpdatedAt,
			&i.Exam.LiveStatus,
			&i.Exam.IpRangeStart,
			&i.Exam.IpRangeEnd,
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

const updateExamLiveStatus = `-- name: UpdateExamLiveStatus :exec
UPDATE exams
SET live_status = $1
WHERE id = $2
`

type UpdateExamLiveStatusParams struct {
	LiveStatus pgtype.Text
	ID         uuid.UUID
}

func (q *Queries) UpdateExamLiveStatus(ctx context.Context, arg UpdateExamLiveStatusParams) error {
	_, err := q.db.Exec(ctx, updateExamLiveStatus, arg.LiveStatus, arg.ID)
	return err
}
