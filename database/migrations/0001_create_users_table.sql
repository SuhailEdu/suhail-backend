-- +goose Up
CREATE TABLE users
(
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name        TEXT          NOT NULL,
    last_name         TEXT          NOT NULL,
    email             citext UNIQUE NOT NULL,
    password          bytea            DEFAULT NULL,
    email_verified_at TIMESTAMP        DEFAULT NULL,
    created_at        TIMESTAMP     NOT NULL,
    updated_at        TIMESTAMP     NOT NULL
);


-- +goose Down

DROP TABLE IF EXISTS users cascade;