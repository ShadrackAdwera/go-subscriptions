package db

import "context"

type CreateUserInput struct {
	UserArgs    CreateUserParams `json:"user"`
	AfterCreate func(user User) error
}

type CreateUserOutput struct {
	Message string `json:"message"`
}

func (s *Store) CreateUserTx(ctx context.Context, args CreateUserInput) (CreateUserOutput, error) {
	err := s.execTx(ctx, func(q *Queries) error {
		user, err := q.CreateUser(ctx, args.UserArgs)

		if err != nil {
			return err
		}

		// send email to verify account + emit User:Create event
		err = args.AfterCreate(user)

		if err != nil {
			return err
		}
		return nil
	})

	return CreateUserOutput{
		Message: "user successfully created",
	}, err
}
