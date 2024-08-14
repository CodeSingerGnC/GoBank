package api

import (
	"fmt"

	db "github.com/CodeSingerGnC/MicroBank/db/sqlc"
	"github.com/CodeSingerGnC/MicroBank/token"
	"github.com/CodeSingerGnC/MicroBank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server 为银行服务提供 Http 请求
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer 创建新的 Http 服务并建立路由
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurreny)
	}

	server.setupRouter()

	return server, nil
}

// setupRouter 设置服务路由
func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/tokens/renew_access", server.renewAccessToken)

	authRouters := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRouters.POST("/accounts", server.createAccount)
	authRouters.GET("/accounts/:id", server.getAccount)
	authRouters.GET("/accounts", server.listAccounts)

	authRouters.POST("/transfers", server.createTransfer)

	server.router = router
}

// Start 在绑定地址上运行 Http 服务
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorResponse 用于向客户端返回错误
func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
