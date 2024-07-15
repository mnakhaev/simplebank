package worker

import (
	"context"

	"github.com/hibiken/asynq"
	db "github.com/mnakhaev/simplebank/db/sqlc"
)

// processor will pick tasks from the Redis queue and process them.

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

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
		asynq.Config{Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  5,
		}},
	)
	return &RedisTaskProcessor{server: server, store: store}
}
