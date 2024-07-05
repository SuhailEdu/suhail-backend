-- +goose Up
ALTER TABLE IF EXISTS exams
    ADD ip_range_start text DEFAULT NULL,
    ADD ip_range_end   text DEFAULT NULL;

-- +goose Down
ALTER TABLE IF EXISTS exams
    DROP COLUMN ip_range_start,
    DROP COLUMN ip_range_end;