package gapi

import (
	"fmt"

	db "github.com/CodeSingerGnC/MicroBank/db/sqlc"
	"github.com/CodeSingerGnC/MicroBank/pb"
	"github.com/CodeSingerGnC/MicroBank/token"
	"github.com/CodeSingerGnC/MicroBank/util"
	"github.com/CodeSingerGnC/MicroBank/worker"
)

// Server 为银行服务提供 grpc 请求
type Server struct {
	pb.UnimplementedMicroBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	taskDistributor worker.TaskDistributor
}

// NewServer 创建新的 grpc 服务
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	// tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
