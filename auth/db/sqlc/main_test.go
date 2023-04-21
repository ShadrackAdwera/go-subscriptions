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
	url := "postgresql://postgres:password@localhost:5432/postgres?sslmode=disable"

	testDb, err = sql.Open("postgres", url)

	testQuery = NewStore(testDb)

	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
