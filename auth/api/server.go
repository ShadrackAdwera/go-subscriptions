package api

import (
	db "github.com/ShadrackAdwera/go-subscriptions/db/sqlc"
	"github.com/ShadrackAdwera/go-subscriptions/token"
	"github.com/ShadrackAdwera/go-subscriptions/workers"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router     *gin.Engine
	store      db.TxStore
	tokenMaker token.TokenMaker
	distro     workers.Distributor
}

func NewServer(store db.TxStore, maker token.TokenMaker, distro workers.Distributor) *Server {
	router := gin.Default()

	srv := &Server{
		store:      store,
		tokenMaker: maker,
		distro:     distro,
	}
	// add routes
	router.POST("/auth/sign-up", srv.signUp)
	router.POST("/auth/login", srv.login)

	srv.router = router
	return srv
}

func errJSON(err error) gin.H {
	return gin.H{"message": err.Error()}
}

func (s *Server) StartServer(addr string) error {
	return s.router.Run(addr)
}
