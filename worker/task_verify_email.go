package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	TaskSendVerifyEmail = "task:send_verify_email"
)

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (td *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marhsal payload: %w", err)
	}
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := td.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("enqueued task")
	return nil
}

func (tp *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		// if body cannot be unmarshalled, then no reason to retry.
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := tp.store.GetUser(ctx, payload.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// if user doesn't exist, then no reason to retry.
			return fmt.Errorf("user doesn't exist: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// TODO: send email to user
	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", user.Email).
		Msg("processed task")
	return nil
}

func (tp *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, tp.ProcessTaskSendVerifyEmail)
	return tp.server.Start(mux)
}
