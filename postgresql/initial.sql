DROP TABLE IF EXISTS test_to_profession;
DROP TABLE IF EXISTS test_answers;
DROP TABLE IF EXISTS test_questions;
DROP TABLE IF EXISTS favourite_professions;
DROP TABLE IF EXISTS professions;
DROP TABLE IF EXISTS users;

CREATE TABLE users
(
    id           SERIAL PRIMARY KEY,
    login        VARCHAR   NOT NULL,
    last_name    VARCHAR   NOT NULL,
    first_name   VARCHAR   NOT NULL,
    patronymic   VARCHAR   NOT NULL,
    step         INTEGER   NOT NULL,
    school_class INTEGER,
    school_name  VARCHAR,
    created_at   TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (login)
);

CREATE TABLE professions
(
    id                SERIAL PRIMARY KEY,
    name              VARCHAR NOT NULL,
    description       VARCHAR NOT NULL,
    short_description VARCHAR,
    tasks             VARCHAR NOT NULL,
    required_skills   VARCHAR NOT NULL,
    relevance         VARCHAR
);

CREATE TABLE favourite_professions
(
    id            SERIAL PRIMARY KEY,
    profession_id INTEGER NOT NULL REFERENCES professions (id),
    user_id       INTEGER NOT NULL REFERENCES users (id),
    UNIQUE (profession_id, user_id)
);

CREATE TABLE test_questions
(
    id      SERIAL PRIMARY KEY,
    test_id INTEGER NOT NULL,
    text    VARCHAR NOT NULL
);

CREATE TABLE test_answers
(
    id          SERIAL PRIMARY KEY,
    user_id     INTEGER NOT NULL REFERENCES users (id),
    question_id INTEGER NOT NULL REFERENCES test_questions (id),
    answer      INTEGER NOT NULL,
    UNIQUE (user_id, question_id)
);

CREATE TABLE test_to_profession
(
    id            SERIAL PRIMARY KEY,
    question_id   INTEGER NOT NULL REFERENCES test_questions (id),
    profession_id INTEGER NOT NULL REFERENCES professions (id),
    correlation   REAL    NOT NULL,
    UNIQUE (question_id, profession_id)
);
