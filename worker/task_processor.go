package worker

import (
	"context"

	db "github.com/CodeSingerGnC/MicroBank/db/sqlc"
	"github.com/CodeSingerGnC/MicroBank/mail"
	"github.com/CodeSingerGnC/MicroBank/util"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	QueueCritical = "critical"
	QueueDefault = "default"
)

type TaskProcessor interface {
	Start() error
	Shutdown()
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server 	*asynq.Server
	store 	db.Store
	mailer 	mail.EmailSender
	config  util.Config
}

func NewRedisTaskProcessor(
	redisOpt asynq.RedisClientOpt, 
	store db.Store, 
	mailer mail.EmailSender, 
	config util.Config,
) TaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault: 5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Info().Err(err).
					Str("type", task.Type()).
					Msg("process task failed")
			}),
			Logger: NewLogger(),
		},
	)

	return &RedisTaskProcessor{
		server: server,
		store:  store,
		mailer: mailer,
		config: config,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Start(mux)
}

func (processor *RedisTaskProcessor) Shutdown() {
	processor.server.Shutdown()
}