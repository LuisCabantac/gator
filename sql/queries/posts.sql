-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, published_at, title, url, description, feed_id)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.*
FROM posts
INNER JOIN feeds ON posts.feed_id = feeds.id
WHERE feeds.user_id = $1
LIMIT $2;
