package workers

import (
	"context"

	db "github.com/ShadrackAdwera/go-subscriptions/db/sqlc"
	"github.com/hibiken/asynq"
)

type Distributor interface {
	DistributeUser(ctx context.Context, user db.User, opts ...asynq.Option) error
}

type TaskDistributor struct {
	client *asynq.Client
}

func NewDistributor(opts *asynq.RedisClientOpt) Distributor {
	client := asynq.NewClient(opts)

	return &TaskDistributor{
		client,
	}
}
