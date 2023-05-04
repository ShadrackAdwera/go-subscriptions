-- name: CreateUserPackage :one
INSERT INTO users_packages (
  user_id, package_id, status, start_date, end_date
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetUserPackages :many
SELECT * 
FROM users_packages 
JOIN users ON users_packages.user_id = users.id
JOIN packages ON users_packages.package_id = packages.id;

-- name: UpdateUserPackage :one
UPDATE users_packages 
SET
  package_id = COALESCE(sqlc.narg(package_id),package_id),
  status = COALESCE(sqlc.narg(status),status),
  start_date = COALESCE(sqlc.narg(start_date),start_date),
  end_date = COALESCE(sqlc.narg(end_date),end_date)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteUserPackage :exec
DELETE FROM users_packages WHERE id = $1;