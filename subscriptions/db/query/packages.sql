-- name: CreatePackage :one
INSERT INTO packages (
  id, name, description, price, stripe_price_id
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetPackages :many
SELECT * FROM packages;

-- name: GetPackage :one
SELECT * 
FROM packages
WHERE id = $1;

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