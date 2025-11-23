-- +goose Up
CREATE TABLE IF NOT EXISTS questions (
    id SERIAL PRIMARY KEY,
    text TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS questions;
