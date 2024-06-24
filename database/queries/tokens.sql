-- name: CreateUserToken :one
INSERT INTO tokens(hash , user_id , expiry , scope )
VALUES ( $1 , $2 , $3 , $4 )
RETURNING *;


-- name: GetUserToken :one
SELECT * FROM tokens WHERE hash = $1 ;
