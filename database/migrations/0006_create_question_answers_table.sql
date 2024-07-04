-- +goose Up
CREATE TABLE IF NOT EXISTS question_answers
(
    user_id     UUId NOT NULL REFERENCES users ON DELETE CASCADE,
    question_id UUId NOT NULL REFERENCES exam_questions ON DELETE CASCADE,
    answer      text NOT NULL,
    is_correct  boolean DEFAULT false,


    PRIMARY KEY (user_id, question_id)
);


-- +goose Down

DROP TABLE IF EXISTS question_answers CASCADE;
