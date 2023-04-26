-- name: CreateUser :one
INSERT INTO users (
  id, username, email
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetUsers :many
SELECT * FROM users;

-- name: UpdateUser :one
UPDATE users 
SET
  username = COALESCE(sqlc.narg(username),username),
  email = COALESCE(sqlc.narg(email),email)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;