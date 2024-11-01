-- name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, last_fetched_at, name, url, user_id)
values (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING *;

-- name: GetFeed :one
SELECT * 
FROM feeds
WHERE name = $1;

-- name: GetFeedByUrl :one
SELECT *
FROM feeds
WHERE URL = $1;

-- name: Pprint :many
SELECT feeds.name, feeds.url, users.name as username
FROM feeds
INNER JOIN users
ON feeds.user_id = users.id;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_at = NOW(), last_fetched_at = NOW()
WHERE feeds.id = $1;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER by last_fetched_at ASC NULLS FIRST
LIMIT 1;