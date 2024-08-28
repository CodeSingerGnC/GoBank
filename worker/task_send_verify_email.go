package worker

import (
	"context"
	"encoding/json"
	"fmt"

	db "github.com/CodeSingerGnC/MicroBank/db/sqlc"
	"github.com/CodeSingerGnC/MicroBank/otpcode"
	"github.com/hibiken/asynq"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/rs/zerolog/log"
)

type PalyloadSendVerifyEmail struct {
	Email string `json:"user_email"`
}

const TaskSendVerifyEmail = "task:send_verify_email"

func generateSecret(issuer, email string) (*otp.Key, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: email,
	})
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PalyloadSendVerifyEmail,
	opt ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opt...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("faided to enqueue task: %w", err)
	}

	log.Info().Str("type: ", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queque", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("enqueue task")
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PalyloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}
	key, err := generateSecret(processor.config.Issuer, payload.Email)
	if err != nil {
		return fmt.Errorf("cannot generate secret key: %w", err)
	}

	passcode, err := otpcode.GeneratePassCode(key.Secret())
	if err != nil {
		return fmt.Errorf("failed to create passcode: %w", err)
	}

	processor.store.DeleteOtpsecret(ctx, payload.Email)
	_, err = processor.store.CreateOtpsecret(ctx, db.CreateOtpsecretParams{
		Email:  payload.Email,
		Secret: key.Secret(),
	})
	if err != nil {
		return fmt.Errorf("failed to store Otpsecret: %w", err)
	}

	subject := "MicroBank Passcode"
	content := fmt.Sprintf(`
	<h1> 请验证您的邮箱 </h1>
	<p>请使用此验证码注册账号 %s <p>
	<p>验证码有效时间为 1 分钟，请尽快完成验证操作。<p>
	`, passcode)

	to := []string{payload.Email}
	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send email")
	}
	log.Info().
		Str("type: ", task.Type()).
		Bytes("payload", task.Payload()).
		Msg("processed task")

	return nil
}
