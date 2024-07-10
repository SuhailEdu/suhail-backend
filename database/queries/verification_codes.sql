-- name: StoreVerificationCode :one
INSERT INTO verification_codes(hash, user_id, expiry, scope)
VALUES ($1, $2, $3, $4)
RETURNING *;


-- name: CheckCodeValidity :one
SELECT EXISTS(SELECT 1 FROM verification_codes WHERE hash = $1 AND expire > now());


-- name: DeleteVerificationCode :exec
DELETE
FROM verification_codes
WHERE hash = $1;
