package gapi

import (
	"context"
	"time"

	db "github.com/CodeSingerGnC/MicroBank/db/sqlc"
	"github.com/CodeSingerGnC/MicroBank/pb"
	"github.com/CodeSingerGnC/MicroBank/util"
	"github.com/CodeSingerGnC/MicroBank/val"
	"github.com/CodeSingerGnC/MicroBank/worker"
	"github.com/go-sql-driver/mysql"
	"github.com/hibiken/asynq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations :=  validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	arg := db.CreateUserTxParam{
		CreateUserParams: db.CreateUserParams{
			UserAccount: req.GetUserAccount(),
			HashPassword: hashedPassword,
			Username: req.GetUsername(),
			Email: req.GetEmail(),
		},
		AfterCreate: func () error{
			// -TODO: use db transaction
			taskPayload := &worker.PalyloadSendVerifyEmail{
				UserAccount: req.GetEmail(),
			}
			opt := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}
			return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opt...)
		},
	}

	err = server.store.CreateUserTx(ctx, arg)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062:
				return nil, status.Errorf(codes.AlreadyExists, "useraccount already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	rsp := &pb.CreateUserResponse{
		UserAccount: req.GetUserAccount(),
	}

	return rsp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUserAccount(req.GetUserAccount()); err != nil {
		violations = append(violations, fieldViolation("user_account", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return 
}
