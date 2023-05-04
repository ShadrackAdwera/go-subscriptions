package api

import (
	"net/http"

	db "github.com/ShadrackAdwera/go-subscriptions/subscriptions/db/sqlc"
	"github.com/gin-gonic/gin"
)

func (srv *Server) getPackages(ctx *gin.Context) {
	packages, err := srv.store.GetPackages(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errJSON(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"packages": packages})
}

type PackageArgs struct {
	Name        string `json:"name" binding:"required,min=5"`
	Description string `json:"description" binding:"required,min=10"`
	Price       int64  `json:"price" binding:"required,min=1"`
}

func (srv *Server) createPackage(ctx *gin.Context) {
	var packageArgs PackageArgs

	if err := ctx.ShouldBindJSON(&packageArgs); err != nil {
		ctx.JSON(http.StatusBadRequest, errJSON(err))
		return
	}

	p, err := srv.store.CreatePackage(ctx, db.CreatePackageParams(packageArgs))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errJSON(err))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"package": p})
}
