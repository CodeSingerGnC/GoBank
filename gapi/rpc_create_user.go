package gapi

import (
	"context"
	"errors"

	db "github.com/CodeSingerGnC/MicroBank/db/sqlc"
	"github.com/CodeSingerGnC/MicroBank/otpcode"
	"github.com/CodeSingerGnC/MicroBank/pb"
	"github.com/CodeSingerGnC/MicroBank/util"
	"github.com/CodeSingerGnC/MicroBank/val"
	"github.com/go-sql-driver/mysql"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	otpsecret, err := server.store.GetOtpsecret(ctx, req.GetEmail())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get otpsecret")
	}

	err = otpcode.VerifyPassCode(req.GetPasscode(), otpsecret.Secret)
	if err != nil {
		if errors.Is(otpcode.ErrPassCodeMismatch, err) {
			server.store.AddOtpsecretTryTime(ctx, req.Email)
			return nil, status.Error(codes.Unauthenticated, "wrong passcode")
		}
		return nil, status.Error(codes.Internal, "failed to verify the passcode")
	}

	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	_, err = server.store.CreateUser(ctx, db.CreateUserParams{
		UserAccount:  req.GetUserAccount(),
		Username:     req.GetUsername(),
		HashPassword: hashedPassword,
		Email:        req.GetEmail(),
	})

	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062:
				return nil, status.Errorf(codes.AlreadyExists, "useraccount already exists: %W", err)
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

	if err := val.ValidatePasscode(req.GetPasscode()); err != nil {
		violations = append(violations, fieldViolation("passcode", err))
	}
	return
}
