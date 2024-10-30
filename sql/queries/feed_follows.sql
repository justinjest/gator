-- name: CreateFeedFollow :many
WITH inserted_feed_follow AS (
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
values ($1, $2, $3, $4, $5)
RETURNING *
)
SELECT inserted_feed_follow.*,
feeds.name AS feed_name,
users.name as user_name
FROM inserted_feed_follow
INNER JOIN feeds on inserted_feed_follow.feed_id = feeds.id
INNER JOIN users on inserted_feed_follow.user_id = users.id;

-- name: GetFeedFollowsForUser :many
SELECT feeds.name, users.name
FROM feed_follows
INNER join users on feed_follows.user_id = users.id
INNER JOIN feeds on feed_follows.feed_id = feeds.id
WHERE users.id = $1;