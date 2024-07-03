-- name: UpdateAnswer :exec
INSERT INTO question_answers (question_id, user_id, answer)
VALUES ($1, $2, $3)
ON CONFLICT (question_id , user_id) DO UPDATE SET answer = $3;