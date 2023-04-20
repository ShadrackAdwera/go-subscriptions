package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/ShadrackAdwera/go-subscriptions/db/mocks"
	db "github.com/ShadrackAdwera/go-subscriptions/db/sqlc"
	"github.com/ShadrackAdwera/go-subscriptions/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func RandomUser(t *testing.T) (db.User, string) {
	username := utils.RandomString(8)
	email := fmt.Sprintf("%s@mail.com", username)
	password := utils.RandomString(12)

	hashedPw, err := utils.HashPassword(password)

	require.NoError(t, err)

	return db.User{
		ID:                utils.RandomInteger(1, 200),
		Username:          username,
		Password:          hashedPw,
		Email:             email,
		PasswordChangedAt: time.Now(),
		CreatedAt:         time.Now(),
	}, password
}

func TestSignUpEndpoint(t *testing.T) {
	user, password := RandomUser(t)

	testCases := []struct {
		name       string
		body       SignUpArgs
		buildStubs func(store *mockdb.MockTxStore)
		comparator func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "TestOK",
			body: SignUpArgs{
				Username: user.Username,
				Email:    user.Email,
				Password: password,
			},
			buildStubs: func(store *mockdb.MockTxStore) {
				// userInput := db.CreateUserInput{
				// 	UserArgs: db.CreateUserParams{
				// 		Username: user.Username,
				// 		Email:    user.Email,
				// 		Password: user.Password,
				// 	},
				// 	AfterCreate: func() error {
				// 		return nil
				// 	},
				// }
				message := db.CreateUserOutput{
					Message: "user successfully created",
				}
				store.EXPECT().FindUserByEmail(gomock.Any(), gomock.Eq(user.Email)).Times(1).Return(db.User{}, nil)
				// rework this . .
				store.EXPECT().CreateUserTx(gomock.Any(), gomock.Any()).Times(1).Return(message, nil)
			},
			comparator: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctlr := gomock.NewController(t)

			store := mockdb.NewMockTxStore(ctlr)
			testCase.buildStubs(store)
			defer ctlr.Finish()

			jsonPayload, err := json.Marshal(testCase.body)

			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/auth/sign-up", bytes.NewReader(jsonPayload))

			require.NoError(t, err)

			recorder := httptest.NewRecorder()

			srv := newServer(store)
			srv.router.ServeHTTP(recorder, req)
			testCase.comparator(t, recorder)
		})
	}
}

func TestLoginEndpoint(t *testing.T) {
	user, password := RandomUser(t)

	testCases := []struct {
		name       string
		body       LoginRequest
		buildStubs func(store *mockdb.MockTxStore)
		comparator func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "TestOK",
			body: LoginRequest{
				Email:    user.Email,
				Password: password,
			},
			buildStubs: func(store *mockdb.MockTxStore) {
				store.EXPECT().FindUserByEmail(gomock.Any(), gomock.Eq(user.Email)).Times(1).Return(user, nil)
				store.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(1)
			},
			comparator: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				compareResponse(t, recorder.Body, user)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctlr := gomock.NewController(t)

			store := mockdb.NewMockTxStore(ctlr)

			testCase.buildStubs(store)
			defer ctlr.Finish()

			jsonPayload, err := json.Marshal(testCase.body)

			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(jsonPayload))

			require.NoError(t, err)

			recorder := httptest.NewRecorder()

			srv := newServer(store)

			srv.router.ServeHTTP(recorder, req)

			testCase.comparator(t, recorder)
		})
	}
}

func compareResponse(t *testing.T, body *bytes.Buffer, user db.User) {
	var response AuthResponse

	jsonRes, err := io.ReadAll(body)

	require.NoError(t, err)

	err = json.Unmarshal(jsonRes, &response)

	require.NoError(t, err)

	require.Equal(t, response.User.ID, user.ID)
	require.Equal(t, response.User.Email, user.Email)
}
