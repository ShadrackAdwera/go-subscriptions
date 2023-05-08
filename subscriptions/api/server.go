package api

import (
	"fmt"

	db "github.com/ShadrackAdwera/go-subscriptions/subscriptions/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	store  db.TxSubscriptionsStore
}

func NewServer(store db.TxSubscriptionsStore) *Server {
	router := gin.Default()

	srv := Server{
		store: store,
	}

	//add middleware
	router.GET("/subscriptions", srv.getPackages)
	router.POST("/subscriptions", srv.createPackage)
	router.POST("/subscriptions/subscribe", srv.subscribePackage)

	srv.router = router
	return &srv
}

func (s *Server) StartServer(addr string) error {
	return s.router.Run(addr)
}

func errJSON(err error) gin.H {
	return gin.H{"message": fmt.Errorf(err.Error())}
}
