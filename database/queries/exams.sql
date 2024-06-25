-- name: CreateExam :one
INSERT INTO exams(id , user_id , title , slug , visibility_status , is_accessable,  created_at , updated_at)
VALUES (uuid_generate_v4() , $1 , $2 , $3 , $4 ,$5 ,  current_timestamp , current_timestamp)
RETURNING *;

-- name: CreateExamQuestions :copyfrom
INSERT INTO exam_questions(  exam_id , question , type , answers   )
VALUES ( $1 , $2 , $3 , $4   )
;

-- name: GetExamById :one
SELECT * FROM exams WHERE id = $1;

-- name: GetUserExams :many
SELECT * FROM exams WHERE user_id = $1;


-- name: GetExamQuestions :many
SELECT * FROM exam_questions WHERE exam_id = $1;

-- name: CheckExamTitleExists :one
SELECT EXISTS(SELECT 1 FROM exams WHERE title = $1 AND user_id = $2);
