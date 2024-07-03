-- +goose Up
ALTER TABLE IF EXISTS exam_participants
    ADD attendance_status text DEFAULT NULL;

-- +goose Down
ALTER TABLE IF EXISTS exam_participants
    DROP COLUMN attendance_status;
