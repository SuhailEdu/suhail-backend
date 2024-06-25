-- +goose Up
CREATE TABLE IF NOT EXISTS tokens (
  hash bytea PRIMARY KEY,
  user_id UUId NOT NULL REFERENCES users ON DELETE CASCADE,
  expiry timestamp NOT NULL,
  scope text NOT NULL
);


-- +goose Down

DROP TABLE IF EXISTS tokens;