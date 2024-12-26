-- name: CreateFollowFeed :one
INSERT INTO follow_feed (id, created_at, updated_at, user_id, feed_id)
VALUES (
	$1,
	$2,
	$3,
	$4,
	$5
)
RETURNING *;

-- name: GetFollowFeed :many
SELECT * FROM follow_feed WHERE user_id=$1;

-- name: DeleteFollowFeed :exec
DELETE FROM follow_feed WHERE id  = $1 AND user_id = $2;
