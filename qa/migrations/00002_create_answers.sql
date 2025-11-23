-- +goose Up
CREATE TABLE IF NOT EXISTS answers (
    id SERIAL PRIMARY KEY,
    question_id INT NOT NULL REFERENCES questions(id),
    user_id TEXT NOT NULL,
    text TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS answers;


