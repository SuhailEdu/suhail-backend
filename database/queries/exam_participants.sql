-- name: CreateExamParticipant :exec
INSERT INTO exam_participants (exam_id, email, status)
VALUES ($1, $2, $3)
ON CONFLICT (email , exam_id) DO NOTHING;

-- name: DeleteParticipants :exec
DELETE
FROM exam_participants
WHERE email = ANY (sqlc.slice(emails))
  AND exam_id = $1
;

-- name: GetExamParticipants :many
SELECT exam_participants.email, exam_participants.status, users.*
FROM exam_participants
         LEFT JOIN users on users.id = exam_participants.user_id
WHERE exam_id = $1
;


-- name: CheckParticipant :one
SELECT EXISTS(SELECT 1 FROM exam_participants WHERE exam_id = $1 AND user_id = $2);

