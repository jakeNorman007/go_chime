-- +goose Up

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY NOT NULL, 
    username VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    password VARCHAR(60) NOT NULL
);

-- +goose Down

DROP TABLE users;
