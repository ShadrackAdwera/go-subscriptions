package api

import (
	"os"
	"testing"

	db "github.com/ShadrackAdwera/go-subscriptions/subscriptions/db/sqlc"
	"github.com/gin-gonic/gin"
)

func newTestServer(store db.TxSubscriptionsStore) *Server {
	srv := NewServer(store)
	return srv
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
