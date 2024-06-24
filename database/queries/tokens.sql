-- name: CreateUserToken :one
INSERT INTO tokens(hash , user_id , expiry , scope )
VALUES ( $1 , $2 , $3 , $4 )
RETURNING *;


-- name: GetUserToken :one
SELECT * FROM tokens WHERE hash = $1 ;

-- name: GetUserByToken :one
SELECT users.id,tokens.hash , tokens.expiry  FROM tokens INNER JOIN users ON users.id = tokens.user_id WHERE hash = $1;
