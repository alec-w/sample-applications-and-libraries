CREATE TABLE posts (
    id integer generated always as identity PRIMARY KEY,
    title varchar NOT NULL,
    content varchar NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);
