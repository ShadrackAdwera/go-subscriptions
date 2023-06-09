// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0
// source: users.sql

package db

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  username, email, password
) VALUES (
  $1, $2, $3
)
RETURNING id, username, password, email, password_changed_at, created_at
`

type CreateUserParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Username, arg.Email, arg.Password)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const findUserByEmail = `-- name: FindUserByEmail :one
SELECT id, username, password, email, password_changed_at, created_at FROM users
WHERE email = $1 LIMIT 1
`

func (q *Queries) FindUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, findUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, username, password, email, password_changed_at, created_at FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users 
SET
  username = COALESCE($1,username),
  email = COALESCE($2,email),
  password = COALESCE($3,password),
  password_changed_at = COALESCE($4,password_changed_at)
WHERE id = $5
RETURNING id, username, password, email, password_changed_at, created_at
`

type UpdateUserParams struct {
	Username          sql.NullString `json:"username"`
	Email             sql.NullString `json:"email"`
	Password          sql.NullString `json:"password"`
	PasswordChangedAt sql.NullTime   `json:"password_changed_at"`
	ID                int64          `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.Username,
		arg.Email,
		arg.Password,
		arg.PasswordChangedAt,
		arg.ID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}
