package worker

import (
	"context"

	"github.com/hibiken/asynq"
	db "github.com/mnakhaev/simplebank/db/sqlc"
)

// processor will pick tasks from the Redis queue and process them.

type TaskProcessor interface {
	// Start is needed to register the task processor to the server.
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}
type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{},
	)
	return &RedisTaskProcessor{server: server, store: store}
}
