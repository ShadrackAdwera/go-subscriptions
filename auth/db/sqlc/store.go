package db

import (
	"context"
	"database/sql"
)

type TxStore interface {
	Querier
	CreateUserTx(ctx context.Context, args CreateUserInput) (CreateUserOutput, error)
}

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) TxStore {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})

	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)

	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}
