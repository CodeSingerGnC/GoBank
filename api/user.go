package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/CodeSingerGnC/MicroBank/db/sqlc"
	"github.com/CodeSingerGnC/MicroBank/util"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type createUserRequest struct {
	UserAccount 	string `json:"user_account" binding:"required,alphanum"`
	Password	 	string `json:"password" binding:"required,min=6"`
	Username		string `json:"username" binding:"required"`
	Email			string `json:"email" binding:"required,email"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		UserAccount: req.UserAccount,
		HashPassword: hashedPassword,
		Username: req.Username,
		Email: req.Email,
	}

	result, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062:
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	_, err = result.LastInsertId()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	ctx.JSON(http.StatusOK, req.UserAccount)
}

type loginUserRequest struct {
	UserAccount 	string `json:"user_account" binding:"required,alphanum"`
	Password	 	string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	SessionID				uuid.UUID	`json:"session_id"`
	AccessToken 		 	string		`json:"access_token"`
	AccessTokenExpiresAt 	time.Time 	`json:"access_token_expires_at"`
	RefreshToken 		 	string		`json:"refresh_token"`
	RefreshTokenExpiresAt 	time.Time 	`json:"refresh_token_expires_at"`
	UserAccount 		 	string		`json:"user_account"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return 
	}

	user, err := server.store.GetUser(ctx, req.UserAccount)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return 
	}

	err = util.CheckPasswordHash(req.Password, user.HashPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.UserAccount,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.UserAccount,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.store.CreateSession(ctx, db.CreateSessionParams{
		ID: refreshPayload.ID[:], 
		UserAccount: user.UserAccount,
		RefreshToken: refreshToken,
		UserAgent: ctx.Request.UserAgent(), 
		ClientIp: ctx.ClientIP(),
		IsBlocked: false, 
		ExpiresAt: refreshPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	sessionID, err := util.BytesToUUID(refreshPayload.ID[:])
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rep := loginUserResponse{
		SessionID: sessionID,
		AccessToken: accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
		RefreshToken: refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt, 
		UserAccount: req.UserAccount,
	}

	ctx.JSON(http.StatusOK, rep)
}