-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY NOT NULL, 
    email TEXT NOT NULL,
    password TEXT NOT NULL
);

-- +goose Down
DROP TABLE users;
