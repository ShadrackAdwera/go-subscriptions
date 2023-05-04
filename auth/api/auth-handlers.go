package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	db "github.com/ShadrackAdwera/go-subscriptions/db/sqlc"
	"github.com/ShadrackAdwera/go-subscriptions/utils"
	"github.com/ShadrackAdwera/go-subscriptions/workers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type SignUpArgs struct {
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (s *Server) signUp(ctx *gin.Context) {
	var signUpArgs SignUpArgs

	if err := ctx.ShouldBindJSON(&signUpArgs); err != nil {
		ctx.JSON(http.StatusBadRequest, errJSON(err))
		return
	}
	//check if email exists
	foundUser, err := s.store.FindUserByEmail(ctx, signUpArgs.Email)

	if err != nil && err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, errJSON(err))
		return
	}

	if foundUser.ID != 0 {
		ctx.JSON(http.StatusForbidden, errJSON(fmt.Errorf("this account already exists")))
		return
	}

	hashedPw, err := utils.HashPassword(signUpArgs.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errJSON(err))
		return
	}

	user, err := s.store.CreateUserTx(ctx, db.CreateUserInput{
		UserArgs: db.CreateUserParams{
			Username: signUpArgs.Username,
			Email:    signUpArgs.Email,
			Password: hashedPw,
		},
		AfterCreate: func(user db.User) error {
			return s.distro.DistributeUser(ctx, workers.UserPayload{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
			}, asynq.MaxRetry(10))
		},
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errJSON(err))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": user.Message})
}

type AuthResponse struct {
	User                 User      `json:"user,omitempty"`
	AccessToken          string    `json:"access_token,omitempty"`
	RefreshToken         string    `json:"refresh_token,omitempty"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at,omitempty"`
	SessionID            uuid.UUID `json:"session_id,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (srv *Server) login(ctx *gin.Context) {
	var authRequest LoginRequest

	if err := ctx.ShouldBindJSON(&authRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errJSON(err))
		return
	}

	user, err := srv.store.FindUserByEmail(ctx, authRequest.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errJSON(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errJSON(err))
		return
	}

	err = utils.IsPassword(authRequest.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errJSON(err))
		return
	}

	aPayload, accessTkn, err := srv.tokenMaker.CreateToken(user.Username, user.ID, user.Email, time.Hour)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errJSON(err))
		return
	}

	pL, refreshTkn, err := srv.tokenMaker.CreateToken(user.Username, user.ID, user.Email, time.Hour*24)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errJSON(err))
		return
	}

	sess, err := srv.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           pL.TokenId,
		Username:     user.Username,
		UserID:       user.ID,
		RefreshToken: refreshTkn,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    pL.ExpiredAt,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errJSON(err))
		return
	}

	loginResponse := AuthResponse{
		User: User{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
		AccessToken:          accessTkn,
		RefreshToken:         refreshTkn,
		AccessTokenExpiresAt: aPayload.ExpiredAt,
		SessionID:            sess.ID,
	}

	ctx.JSON(http.StatusOK, loginResponse)
}
