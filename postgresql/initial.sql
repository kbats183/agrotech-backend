DROP TABLE IF EXISTS favourite_study_programs;
DROP TABLE IF EXISTS test_to_profession;
DROP TABLE IF EXISTS test_answers;
DROP TABLE IF EXISTS test_questions;
DROP TABLE IF EXISTS favourite_professions;
DROP TABLE IF EXISTS specializations_professions;
DROP TABLE IF EXISTS specialization;
DROP TABLE IF EXISTS university;

DROP TABLE IF EXISTS vacancies;
DROP TABLE IF EXISTS professions;
DROP TABLE IF EXISTS users;

CREATE TABLE users
(
    id                       SERIAL PRIMARY KEY,
    login                    VARCHAR   NOT NULL,
    last_name                VARCHAR   NOT NULL,
    first_name               VARCHAR   NOT NULL,
    patronymic               VARCHAR   NOT NULL,
    region                   INTEGER,
    step                     INTEGER,
    school_class             INTEGER,
    school_name              VARCHAR,
    university_name          VARCHAR,
    university_study_program VARCHAR,
    university_profession    INTEGER,
    created_at               TIMESTAMP NOT NULL DEFAULT now(),
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

CREATE TABLE specialization
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    disciplines VARCHAR NOT NULL
);

CREATE TABLE specializations_professions
(
    id                SERIAL PRIMARY KEY,
    specialization_id INTEGER NOT NULL REFERENCES specialization (id),
    professions_id    INTEGER NOT NULL REFERENCES professions (id)
);

CREATE TABLE university
(
    id      SERIAL PRIMARY KEY,
    name    VARCHAR NOT NULL,
    address VARCHAR NOT NULL,
    image   VARCHAR NOT NULL
);

CREATE TABLE study_programs
(
    id                SERIAL PRIMARY KEY,
    specialization_id INTEGER NOT NULL REFERENCES specialization (id),
    university_id     INTEGER NOT NULL REFERENCES university (id),
    exams             JSONB   NOT NULL,
    score_budget      INTEGER,
    score_contract    INTEGER,
    contract_amount   INTEGER,
    grade_budget      REAL,
    grade_contract    REAL
);

CREATE TABLE favourite_study_programs
(
    id                SERIAL PRIMARY KEY,
    study_programs_id INTEGER NOT NULL REFERENCES study_programs (id),
    user_id           INTEGER NOT NULL REFERENCES users (id),
    UNIQUE (study_programs_id, user_id)
);

CREATE TABLE vacancies
(
    id             SERIAL PRIMARY KEY,
    hh_id          INTEGER,
    name           VARCHAR NOT NULL,
    url            VARCHAR NOT NULL,
    employer       VARCHAR NOT NULL,
    employer_url   VARCHAR NOT NULL,
    employer_image VARCHAR NOT NULL,
    responsibility VARCHAR NOT NULL,
    area           INT     NOT NULL,
    UNIQUE (hh_id)
)