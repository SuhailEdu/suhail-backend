-- +goose Up
CREATE TABLE IF NOT EXISTS question_answers
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUId NOT NULL REFERENCES users ON DELETE CASCADE,
    question_id UUId NOT NULL REFERENCES exam_questions ON DELETE CASCADE,
    answer      text NOT NULL,
    is_correct  boolean          DEFAULT false
);


-- +goose Down

DROP TABLE IF EXISTS question_answers CASCADE;
