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
    last_fetched_at timestamp,
    name text NOT NULL,
    url text UNIQUE NOT NULL,
    user_id text NOT NULL references users(id) ON DELETE CASCADE
);

CREATE TABLE feed_follows(
    id text PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    user_id text NOT NULL references users(id) ON DELETE CASCADE,
    feed_id text NOT NULL references feeds(id) ON DELETE CASCADE,
    UNIQUE(user_id, feed_id) 
);

-- +goose Down
DROP TABLE feeds CASCADE;
DROP TABLE feed_follows CASCADE;
DROP TABLE users CASCADE;
