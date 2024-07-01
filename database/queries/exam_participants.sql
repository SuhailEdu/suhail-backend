-- name: CreateExamParticipant :exec
INSERT INTO exam_participants (exam_id, email, status)
VALUES ($1, $2, $3)
ON CONFLICT (email , exam_id) DO NOTHING;

-- name: DeleteParticipants :exec
DELETE
FROM exam_participants
WHERE user_id = ANY (SELECT id FROM users WHERE users.email = ANY (sqlc.slice(emails)))
  AND exam_id = $1
;

-- name: GetExamParticipants :many
SELECT sqlc.embed(users)
FROM exam_participants
         INNER JOIN users on users.id = exam_participants.user_id
WHERE exam_id = $1
;

