package db

import "database/sql"

type TxSubscriptionsStore interface {
	Querier
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
