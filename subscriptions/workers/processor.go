package workers

import (
	"context"

	db "github.com/ShadrackAdwera/go-subscriptions/subscriptions/db/sqlc"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	QueueCritical = "critical"
	QueueDefaullt = "default"
	QueueLow      = "low"
)

type TaskProcessor interface {
	CreateAsyncUser(ctx context.Context, task *asynq.Task) error
	Start() error
}

type Processor struct {
	srv   *asynq.Server
	store db.TxStore
}

func NewProcessor(opts *asynq.RedisClientOpt, store db.TxStore) TaskProcessor {
	server := asynq.NewServer(opts, asynq.Config{
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Err(err).Str("task_type", task.Type()).Bytes("payload", task.Payload()).Msg("error processing task . . ")
		}),
		Queues: map[string]int{
			QueueCritical: 6,
			QueueDefaullt: 3,
			QueueLow:      1,
		},
	})

	return &Processor{
		srv:   server,
		store: store,
	}
}

func (p *Processor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskCreateUser, p.CreateAsyncUser)

	return p.srv.Start(mux)
}
