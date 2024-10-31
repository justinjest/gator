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

-- MarkFeedFetched :exec
UPDATE feeds
SET feeds.updated_at = NOW(), feeds.last_fetched_at = NOW()
WHERE feeds.id = $1;