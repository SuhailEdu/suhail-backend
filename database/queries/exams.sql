-- name: CreateExam :one
INSERT INTO exams(id , user_id , title , slug , visibility_status , is_accessable,  created_at , updated_at)
VALUES (uuid_generate_v4() , $1 , $2 , $3 , $4 ,$5 ,  current_timestamp , current_timestamp)
RETURNING *;

-- name: CreateExamQuestions :one
INSERT INTO exam_questions(id ,  exam_id , question , type , answers ,  created_at , updated_at)
VALUES (uuid_generate_v4() , $1 , $2 , $3 , $4 , current_timestamp , current_timestamp)
RETURNING *;

-- name: GetExamById :one
SELECT * FROM exams WHERE id = $1;

-- name: GetUserExams :many
SELECT * FROM exams WHERE user_id = $1;
