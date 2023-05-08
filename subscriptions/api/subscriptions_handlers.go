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
	//add authorization

	if err := ctx.ShouldBindJSON(&packageArgs); err != nil {
		ctx.JSON(http.StatusBadRequest, errJSON(err))
		return
	}

	p, err := srv.store.CreatePackageTx(ctx, db.CreatePackageTxInput{
		Name:        packageArgs.Name,
		Description: packageArgs.Description,
		Price:       packageArgs.Price,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errJSON(err))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": p.Message})
}

type SubscriptionPackageArgs struct {
	CustomerID          string `json:"customer_id" binding:"required,min=5"`
	SubscriptionPackage string `json:"subscription_package" binding:"required,min=5"`
}

func (srv *Server) subscribePackage(ctx *gin.Context) {
	var subscriptionPackageArgs SubscriptionPackageArgs

	if err := ctx.ShouldBindJSON(&subscriptionPackageArgs); err != nil {
		ctx.JSON(http.StatusBadRequest, errJSON(err))
		return
	}

	response, err := srv.store.SubscribePackageTx(ctx, db.SubscribePackageTxInput{
		CustomerID:            subscriptionPackageArgs.CustomerID,
		SubscriptionPackageID: subscriptionPackageArgs.SubscriptionPackage,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errJSON(err))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": response.Message})
}
