package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

// this file contains code to create tasks and distribute them to workers via Redis queue

type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

// NewRedisTaskDistributor returns an interface in order to make `RedisTaskDistributor` to implement that interface.
func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{client: client}
}
