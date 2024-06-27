-- name: CreateExamParticipant :exec
INSERT INTO exam_participants (exam_id, user_id)
VALUES ($1,
        (SELECT (id)
         FROM users
         WHERE email = $2
         LIMIT 1))
ON CONFLICT (user_id , exam_id) DO NOTHING;
