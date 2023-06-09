// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: packages.sql

package db

import (
	"context"
	"database/sql"
)

const createPackage = `-- name: CreatePackage :one
INSERT INTO packages (
  id, name, description, price, stripe_price_id
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id, name, description, price, stripe_price_id
`

type CreatePackageParams struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Price         int64  `json:"price"`
	StripePriceID string `json:"stripe_price_id"`
}

func (q *Queries) CreatePackage(ctx context.Context, arg CreatePackageParams) (Package, error) {
	row := q.db.QueryRowContext(ctx, createPackage,
		arg.ID,
		arg.Name,
		arg.Description,
		arg.Price,
		arg.StripePriceID,
	)
	var i Package
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.StripePriceID,
	)
	return i, err
}

const deletePackage = `-- name: DeletePackage :exec
DELETE FROM packages WHERE id = $1
`

func (q *Queries) DeletePackage(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deletePackage, id)
	return err
}

const getPackage = `-- name: GetPackage :one
SELECT id, name, description, price, stripe_price_id 
FROM packages
WHERE id = $1
`

func (q *Queries) GetPackage(ctx context.Context, id string) (Package, error) {
	row := q.db.QueryRowContext(ctx, getPackage, id)
	var i Package
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.StripePriceID,
	)
	return i, err
}

const getPackages = `-- name: GetPackages :many
SELECT id, name, description, price, stripe_price_id FROM packages
`

func (q *Queries) GetPackages(ctx context.Context) ([]Package, error) {
	rows, err := q.db.QueryContext(ctx, getPackages)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Package
	for rows.Next() {
		var i Package
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.StripePriceID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePackage = `-- name: UpdatePackage :one
UPDATE packages 
SET
  name = COALESCE($1,name),
  description = COALESCE($2,description),
  price = COALESCE($3,price)
WHERE id = $4
RETURNING id, name, description, price, stripe_price_id
`

type UpdatePackageParams struct {
	Name        sql.NullString `json:"name"`
	Description sql.NullString `json:"description"`
	Price       sql.NullInt64  `json:"price"`
	ID          string         `json:"id"`
}

func (q *Queries) UpdatePackage(ctx context.Context, arg UpdatePackageParams) (Package, error) {
	row := q.db.QueryRowContext(ctx, updatePackage,
		arg.Name,
		arg.Description,
		arg.Price,
		arg.ID,
	)
	var i Package
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.StripePriceID,
	)
	return i, err
}
