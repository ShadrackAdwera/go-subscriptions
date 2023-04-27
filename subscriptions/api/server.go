package api

import (
	"net/http"

	db "github.com/ShadrackAdwera/go-subscriptions/subscriptions/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	store  db.TxStore
}

func NewServer(store db.TxStore) *Server {
	router := gin.Default()

	srv := Server{
		store: store,
	}

	//add middleware
	router.GET("/subscriptions", func(ctx *gin.Context) {
		ctx.JSON(http.StatusAccepted, gin.H{"message": "ping subscription service"})
	})

	srv.router = router
	return &srv
}

func (s *Server) StartServer(addr string) error {
	return s.router.Run(addr)
}
