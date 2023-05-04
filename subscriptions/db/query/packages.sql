-- name: CreatePackage :one
INSERT INTO packages (
  name, description, price
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetPackages :many
SELECT * FROM packages;

-- name: UpdatePackage :one
UPDATE packages 
SET
  name = COALESCE(sqlc.narg(name),name),
  description = COALESCE(sqlc.narg(description),description),
  price = COALESCE(sqlc.narg(price),price)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeletePackage :exec
DELETE FROM packages WHERE id = $1;