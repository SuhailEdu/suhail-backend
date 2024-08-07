-- +goose Up
CREATE TABLE IF NOT EXISTS exams
(
    id                UUID PRIMARY KEY   DEFAULT uuid_generate_v4(),
    user_id           UUId      NOT NULL REFERENCES users ON DELETE CASCADE,
    title             text      NOT NULL,
    slug              text               DEFAULT NULL,
    visibility_status text      NOT NULL,
    is_accessable     boolean            DEFAULT true,

    created_at        TIMESTAMP NOT NULL DEFAULT now(),
    updated_at        TIMESTAMP NOT NULL DEFAULT now()
);


CREATE TABLE IF NOT EXISTS exam_questions
(
    id         UUID PRIMARY KEY   DEFAULT uuid_generate_v4(),
    exam_id    UUId      NOT NULL REFERENCES exams ON DELETE CASCADE,
    question   TEXT      NOT NULL,
    answers    JSON      NOT NULL,
    type       text      NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS exam_participants
(
    user_id    UUId REFERENCES users ON DELETE CASCADE,
    exam_id    UUId      NOT NULL REFERENCES exams ON DELETE CASCADE,
    email      citext    NOT NULL,
    status     text      NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),

    PRIMARY KEY (email, exam_id)
);


-- +goose Down

DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS exam_participants;
DROP TABLE IF EXISTS exam_questions;
