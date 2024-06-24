-- name: CreateUser :one
INSERT INTO users(id , first_name , last_name , password , email ,  created_at , updated_at)
VALUES (uuid_generate_v4() , $1 , $2 , $3 , $4 , current_timestamp , current_timestamp)
RETURNING *;


-- name: CheckEmailUniqueness :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);


-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;
