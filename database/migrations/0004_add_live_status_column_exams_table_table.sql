-- +goose Up
ALTER TABLE exams
    ADD live_status text DEFAULT NULL;

-- +goose Down
ALTER TABLE IF EXISTS exams
    DROP COLUMN live_status;
