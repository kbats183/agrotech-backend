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
DROP TABLE IF EXISTS professions;
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

DROP TABLE IF EXISTS favourite_professions;
CREATE TABLE favourite_professions
(
    id            SERIAL PRIMARY KEY,
    profession_id INTEGER NOT NULL REFERENCES professions (id),
    user_id       INTEGER NOT NULL REFERENCES users (id),
    UNIQUE (profession_id, user_id)
);
