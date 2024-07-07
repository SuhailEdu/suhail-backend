-- name: CreateExam :one
INSERT INTO exams(id, user_id, title, slug, visibility_status, is_accessable, created_at, updated_at)
VALUES (uuid_generate_v4(), $1, $2, $3, $4, $5, current_timestamp, current_timestamp)
RETURNING *;

-- name: GetExamById :one
SELECT *
FROM exams
WHERE id = $1;

-- name: UpdateExam :exec
UPDATE exams
SET title             = $1,
    visibility_status = $2,
    ip_range_start    = $3,
    ip_range_end      = $4
WHERE id = $5

RETURNING *
;


-- name: GetUserExams :many
SELECT exams.*, COUNT(exam_participants.*) as particpants_count, COUNT(exam_questions.*) as questions_count
FROM exams
         LEFT JOIN exam_questions ON exam_questions.exam_id = exams.id
         LEFT JOIN exam_participants ON exam_participants.exam_id = exams.id
WHERE exams.user_id = $1
GROUP BY exams.id
ORDER BY exams.created_at DESC
;

-- name: GetParticipatedExams :many
SELECT exams.*, COUNT(exam_questions.*) as questions_count
FROM exams
         LEFT JOIN exam_questions ON exam_questions.exam_id = exams.id
         INNER JOIN exam_participants ON exam_participants.exam_id = exams.id AND exam_participants.user_id = $1
-- WHERE user_id = $1
GROUP BY exams.id
ORDER BY exams.created_at DESC
;


-- name: GetUserExamsWithQuestions :many
SELECT *
FROM exams
         LEFT JOIN exam_questions ON exam_questions.exam_id = exams.id
WHERE user_id = $1;

-- name: CheckExamTitleExists :one
SELECT EXISTS(SELECT 1
              FROM exams
              WHERE title = $1
                AND user_id = $2
                AND (@id::uuid is null or id != @id::uuid)


--                 AND id != CASE WHEN @id::string THEN @id::string ELSE id END
--                                   AND (case when @except::uuid then where id != @except::uuid)
);


-- name: FindMyExam :one
SELECT exams.*, COUNT(exam_questions.*) as questions_count, COUNT(exam_participants.*) as participants_count
FROM exams
         LEFT JOIN exam_questions ON exam_questions.exam_id = exams.id
         LEFT JOIN exam_participants ON exam_participants.exam_id = exams.id
WHERE exams.id = $1
  AND exams.user_id = $2
GROUP BY exams.id
;

-- name: FindMyParticipatedExam :many
SELECT sqlc.embed(exams), sqlc.embed(exam_questions)
FROM exams
         INNER JOIN exam_participants ON exam_participants.exam_id = exams.id
         LEFT JOIN exam_questions ON exam_questions.exam_id = exams.id
WHERE exam_participants.user_id = $1
  AND exams.id = $2
;


-- name: DeleteExam :exec
DELETE
FROM exams
WHERE id = $1;
