-- +goose Up
CREATE TABLE verification_codes
(
    hash       bytea PRIMARY KEY,
    user_id    UUID      NOT NULL REFERENCES users ON DELETE CASCADE,
    scope      text      NOT NULL,
    expiry     timestamp NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);


-- +goose Down

DROP TABLE IF EXISTS verification_codes CASCADE;
