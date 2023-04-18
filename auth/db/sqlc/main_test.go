package db

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQuery TxStore
var testDb *sql.DB

func TestMain(m *testing.M) {
	var err error
	url := "postgresql://root:password@localhost:5431/go_subscriptions_auth?sslmode=disable"

	testDb, err = sql.Open("postgres", url)

	testQuery = NewStore(testDb)

	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
