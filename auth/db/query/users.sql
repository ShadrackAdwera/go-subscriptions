-- name: CreateUser :one
INSERT INTO users (
  username, email, password
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: FindUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users 
SET
  username = COALESCE(sqlc.narg(username),username),
  email = COALESCE(sqlc.narg(email),email),
  password = COALESCE(sqlc.narg(password),password),
  password_changed_at = COALESCE(sqlc.narg(password_changed_at),password_changed_at)
WHERE id = sqlc.arg(id)
RETURNING *;