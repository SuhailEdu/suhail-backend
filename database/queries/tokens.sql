-- name: CreateUserToken :one
INSERT INTO tokens(hash , user_id , expiry , scope )
VALUES ( $1 , $2 , $3 , $4 )
RETURNING *;


-- name: CheckTokenIsValid :one
SELECT EXISTS(SELECT 1 FROM tokens WHERE hash = $1 AND expiry > now());
