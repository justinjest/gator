-- +goose Up
CREATE TABLE users(
id text PRIMARY KEY,
created_at timestamp NOT NULL,
updated_at timestamp NOT NULL,
name text UNIQUE NOT NULL
);

CREATE TABLE feeds(
    id text PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    name text NOT NULL,
    url text UNIQUE NOT NULL,
    user_id text references users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE users;
