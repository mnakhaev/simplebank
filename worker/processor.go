package worker

import (
	"context"

	"github.com/hibiken/asynq"
	db "github.com/mnakhaev/simplebank/db/sqlc"
	"github.com/rs/zerolog/log"
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
		},
			ErrorHandler: asynq.ErrorHandlerFunc(func(_ context.Context, task *asynq.Task, _ error) {
				log.Error().
					Str("type", task.Type()).
					Bytes("payload", task.Payload()).
					Msg("process task failed")
			}), // handle processor errors
			Logger: NewLogger(), // use custom logger to force asynq to follow needed logger format
		},
	)
	return &RedisTaskProcessor{server: server, store: store}
}
