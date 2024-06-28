-- name: CreateExamParticipant :exec
INSERT INTO exam_participants (exam_id, user_id)
VALUES ($1,
        (SELECT (id)
         FROM users
         WHERE email = $2
         LIMIT 1))
ON CONFLICT (user_id , exam_id) DO NOTHING;

-- name: DeleteParticipants :exec
DELETE
FROM exam_participants
WHERE exam_id = $1
--   AND user_id = ANY ($2::int[])
--   AND user_id IN (SELECT distinct id FROM users WHERE users.email in (sqlc.slice(emails)));
  AND user_id IN (sqlc.slice(emails));
;

