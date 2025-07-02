-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: ListFeeds :many
SELECT f.*, u.name as "user_name" 
FROM feeds f
join users u on f.user_id = u.id
ORDER BY f.created_at DESC;

-- name: GetFeedByUrl :one
SELECT f.*
FROM feeds f
WHERE f.url = $1
LIMIT 1;