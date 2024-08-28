package gapi

import (
	"context"
	"time"

	"github.com/CodeSingerGnC/MicroBank/pb"
	"github.com/CodeSingerGnC/MicroBank/val"
	"github.com/CodeSingerGnC/MicroBank/worker"
	"github.com/hibiken/asynq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *Server) SendPassCode(ctx context.Context, req *pb.SendPassCodeRequest) (*pb.SendPassCodeResponse, error) {
	violations := validateSendPassCodeRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	taskPayload := &worker.PalyloadSendVerifyEmail{
		Email: req.GetEmail(),
	}
	opts := []asynq.Option{
		asynq.MaxRetry(3),
		asynq.ProcessIn(1 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}
	server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)

	return &pb.SendPassCodeResponse{}, nil
}

func validateSendPassCodeRequest(req *pb.SendPassCodeRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("user_email", err))
	}
	return
}
