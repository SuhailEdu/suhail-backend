-- name: UpdateAnswer :exec
INSERT INTO question_answers (question_id, user_id, answer)
VALUES ($1, $2, $3)
ON CONFLICT (question_id , user_id) DO UPDATE SET answer = $3;


-- name: GetParticipantAnswers :many
SELECT exams.id          as exam_id,
       question_answers.user_id,
       exam_questions.id as question_id,
       exam_questions.type,
       exam_questions.answers,
       question_answers.answer,
       exam_questions.question
FROM exams
         LEFT JOIN exam_questions ON exam_questions.exam_id = exams.id
         LEFT JOIN question_answers ON question_answers.question_id = exam_questions.id
WHERE exams.id = $1
  AND question_answers.user_id = $2;


