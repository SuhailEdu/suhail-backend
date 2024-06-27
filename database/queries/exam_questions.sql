-- name: UpdateQuestion :exec
UPDATE exam_questions
SET question = $1,
    answers  = $2
WHERE id = $3

RETURNING *
;

-- name: CreateExamQuestions :copyfrom
INSERT INTO exam_questions(exam_id, question, type, answers)
VALUES ($1, $2, $3, $4)
;


-- name: CreateQuestion :one
INSERT INTO exam_questions(exam_id, question, type, answers)
VALUES ($1, $2, $3, $4)
RETURNING *
;

-- name: GetQuestionById :one
SELECT *
FROM exam_questions
WHERE id = $1
LIMIT 1
;


-- name: CheckQuestionTitleExits :one
SELECT EXISTS(SELECT 1 FROM exam_questions WHERE question = $1 AND exam_id = $2);
;

-- name: DeleteQuestion :exec
DELETE
FROM exam_questions
WHERE id = $1
;
