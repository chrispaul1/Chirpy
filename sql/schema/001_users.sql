-- +goose Up
CREATE TABLE users(
    ID UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    email TEXT UNIQUE NOT NULL
);

ALTER TABLE users 
ADD COLUMN hashed_password TEXT UNIQUE NOT NULL DEFAULT 'unset';

ALTER TABLE users
ADD COLUMN is_chirpy_red BOOL NOT NULL DEFAULT FALSE;

-- +goose Down
DROP TABLE users;
