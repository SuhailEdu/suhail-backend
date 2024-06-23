-- +goose Up
CREATE TABLE USERS(
    id UUID PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email citext UNIQUE NOT NULL,
    password bytea DEFAULT NULL,
    email_verified_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);



-- +goose Down

DROP TABLE USERS;