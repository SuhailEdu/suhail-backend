-- name: UpdateQuestion :exec
UPDATE exam_questions
SET question = $1,
    answers  = $2
WHERE id = $3

RETURNING *
;


-- name: GetQuestionById :one
SELECT *
FROM exam_questions
WHERE id = $1
LIMIT 1
;

-- name: DeleteQuestion :exec
DELETE
FROM exam_questions
WHERE id = $1
;
