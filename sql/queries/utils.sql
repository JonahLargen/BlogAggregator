-- name: ResetAll :exec
TRUNCATE TABLE feeds, users, feed_follows CASCADE;