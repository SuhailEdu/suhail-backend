-- name: GetLiveExamQuestionForManager :many
SELECT sqlc.embed(exam_questions), sqlc.embed(exams)
FROM exam_questions
         INNER JOIN exams ON exams.id = $1
LIMIT 1
;
-- name: GetLiveExamParticipants :many
SELECT sqlc.embed(users), sqlc.embed(exam_participants)
FROM exam_participants
         INNER JOIN users ON users.id = exam_participants.user_id
WHERE exam_participants.exam_id = $1
;
