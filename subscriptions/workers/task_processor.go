package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	db "github.com/ShadrackAdwera/go-subscriptions/subscriptions/db/sqlc"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	stripe "github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
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

	// create customer in stripe
	sKey := os.Getenv("STRIPE_KEY")

	if sKey == "" {
		return fmt.Errorf("no stripe key found %w", asynq.SkipRetry)
	}

	stripe.Key = sKey
	params := &stripe.CustomerParams{
		Name:  stripe.String(payload.Username),
		Email: stripe.String(payload.Email),
	}

	c, err := customer.New(params)

	if err != nil {
		return fmt.Errorf("an error occured while creating the user in stripe %w", err)
	}

	user, err := processor.store.CreateSubscriptionUser(ctx, db.CreateSubscriptionUserParams{
		ID:       payload.ID,
		Username: payload.Username,
		Email:    payload.Email,
		StripeID: c.ID,
	})

	if err != nil {
		return fmt.Errorf("an error occured while creating the user %w", err)
	}

	log.Info().
		Str("user_name", user.Username).
		Int64("with ID of", user.ID).
		Str("with Stripe ID of", c.ID).
		Str("created at ", user.CreatedAt.Format("2006-01-02 15:04:05")).
		Int64("with ID of", user.ID).
		Msg("Async User Created")

	return nil
}
