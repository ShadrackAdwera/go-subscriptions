package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/ShadrackAdwera/go-subscriptions/utils"
	"github.com/stretchr/testify/require"
)

func CreateTestUser(t *testing.T) User {
	username := utils.RandomString(8)
	passw := utils.RandomString(12)

	hashPw, err := utils.HashPassword(passw)

	require.NoError(t, err)

	user, err := testQuery.CreateUser(context.Background(), CreateUserParams{
		Username: username,
		Email:    fmt.Sprintf("%s@mail.com", username),
		Password: hashPw,
	})

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.NotZero(t, user.ID)

	return user
}

func TestCreateUser(t *testing.T) {
	CreateTestUser(t)
}
