package api

import (
	"log"
	"os"
	"testing"

	db "github.com/ShadrackAdwera/go-subscriptions/db/sqlc"
	"github.com/ShadrackAdwera/go-subscriptions/token"
	"github.com/ShadrackAdwera/go-subscriptions/utils"
	"github.com/gin-gonic/gin"
)

func newServer(store db.TxStore) *Server {
	tknMaker, err := token.NewPasetoMaker(utils.RandomString(32))

	if err != nil {
		log.Panic(err)
	}

	srv := NewServer(store, tknMaker, nil)
	return srv
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
