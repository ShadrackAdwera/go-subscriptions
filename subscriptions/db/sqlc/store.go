package db

import (
	"context"
	"database/sql"
)

type TxSubscriptionsStore interface {
	Querier
	CreatePackageTx(ctx context.Context, args CreatePackageTxInput) (CreatePackageTxResult, error)
	SubscribePackageTx(ctx context.Context, args SubscribePackageTxInput) (SubscribePackageTxOutput, error)
}

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) TxSubscriptionsStore {
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
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}
	return tx.Commit()
}
