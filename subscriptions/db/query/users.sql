-- name: CreateSubscriptionUser :one
INSERT INTO users (
  id, username, email
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetSubscriptionUsers :many
SELECT * FROM users;

-- name: UpdateSubscriptionUser :one
UPDATE users 
SET
  username = COALESCE(sqlc.narg(username),username),
  email = COALESCE(sqlc.narg(email),email)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteSubscriptionUser :exec
DELETE FROM users WHERE id = $1;