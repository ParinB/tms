-- +goose Up
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email citext not null unique,
    firstname varchar not null,
    lastname varchar not null,
    created_at timestamp not null default now()::timestamp,
    updated_at timestamp not null default now()::timestamp
);


-- +goose Down
DROP TABLE users;