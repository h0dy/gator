-- name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllFeeds :many
SELECT 
feeds.name AS feed_name,
feeds.url AS feed_url,
feeds.created_at,
feeds.updated_at,
users.name AS created_by FROM feeds
JOIN users ON feeds.user_id = users.id;

-- name: GetFeedByURL :one 
SELECT * FROM feeds WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds 
SET updated_at = NOW(),
    last_fetched_at = NOW() 
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST 
LIMIT 1;