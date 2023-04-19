package api

import (
	db "github.com/ShadrackAdwera/go-subscriptions/db/sqlc"
	"github.com/ShadrackAdwera/go-subscriptions/token"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router     *gin.Engine
	store      db.TxStore
	tokenMaker token.TokenMaker
}

func NewServer(store db.TxStore, maker token.TokenMaker) *Server {
	router := gin.Default()

	srv := &Server{
		store:      store,
		tokenMaker: maker,
	}

	// add routes
	srv.router = router
	return srv
}

func (s *Server) StartServer(addr string) error {
	return s.router.Run(addr)
}
