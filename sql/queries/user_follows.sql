-- name: CreateFeedFollow :one
INSERT INTO user_follows (id, created_at, updated_at, user_id, feed_id) 
VALUES($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFeedFollowsForUser :many
SELECT feeds.name as feedName, users.name as userName FROM user_follows
INNER JOIN feeds ON feeds.id = user_follows.feed_id
INNER JOIN users ON users.id = user_follows.user_id
WHERE users.name = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM user_follows 
WHERE user_follows.user_id = $1 AND feed_id IN (
    SELECT id FROM feeds WHERE url = $2
);
