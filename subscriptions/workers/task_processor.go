package workers

import (
	"context"
	"encoding/json"
	"fmt"

	db "github.com/ShadrackAdwera/go-subscriptions/subscriptions/db/sqlc"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type UserCreatedPayload struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

const TaskCreateUser = "task:create_user"

func (processor *Processor) CreateAsyncUser(ctx context.Context, task *asynq.Task) error {
	var payload UserCreatedPayload

	err := json.Unmarshal(task.Payload(), &payload)

	if err != nil {
		return fmt.Errorf("error unmarshalling JSON %w", asynq.SkipRetry)
	}

	user, err := processor.store.CreateSubscriptionUser(ctx, db.CreateSubscriptionUserParams(payload))

	if err != nil {
		return fmt.Errorf("an error occured while creating the user %w", err)
	}

	log.Info().
		Str("user_name", user.Username).
		Str("email", user.Email).
		Str("created at ", user.CreatedAt.Format("2006-01-02 15:04:05")).
		Int64("with ID of", user.ID).
		Msg("Async User Created")

	return nil
}
