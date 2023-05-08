-- name: CreateSubscriptionUser :one
INSERT INTO users (
  id, username, email, stripe_id
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetSubscriptionUsers :many
SELECT * FROM users;

-- name: GetSubscriptionUserByStripeId :one
SELECT * 
FROM users 
WHERE stripe_id = $1 
LIMIT 1;

-- name: UpdateSubscriptionUser :one
UPDATE users 
SET
  username = COALESCE(sqlc.narg(username),username),
  email = COALESCE(sqlc.narg(email),email),
  stripe_id = COALESCE(sqlc.narg(stripe_id),stripe_id)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteSubscriptionUser :exec
DELETE FROM users WHERE id = $1;