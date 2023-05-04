package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mocksdb "github.com/ShadrackAdwera/go-subscriptions/subscriptions/db/mocks"
	db "github.com/ShadrackAdwera/go-subscriptions/subscriptions/db/sqlc"
	"github.com/ShadrackAdwera/go-subscriptions/subscriptions/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func dummyPackage() db.Package {
	return db.Package{
		ID:          utils.RandomInteger(1, 100),
		Name:        utils.RandomString(10),
		Description: utils.RandomString(20),
		Price:       utils.RandomInteger(5, 20),
	}
}

func TestGetPackages(t *testing.T) {
	n := 3
	packages := make([]db.Package, n)

	for i := 0; i < n; i++ {
		packages[i] = dummyPackage()
	}

	testCases := []struct {
		name       string
		buildStubs func(store *mocksdb.MockTxSubscriptionsStore)
		comparator func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "TESTOK",
			buildStubs: func(store *mocksdb.MockTxSubscriptionsStore) {
				store.EXPECT().GetPackages(gomock.Any()).Times(1).Return(packages, nil)
			},
			comparator: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				compareResponses(t, recorder.Body, packages)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctlr := gomock.NewController(t)
			store := mocksdb.NewMockTxSubscriptionsStore(ctlr)

			testCase.buildStubs(store)

			defer ctlr.Finish()

			srv := newTestServer(store)

			req, err := http.NewRequest(http.MethodGet, "/subscriptions", nil)

			require.NoError(t, err)

			res := httptest.NewRecorder()

			srv.router.ServeHTTP(res, req)

			testCase.comparator(t, res)
		})
	}
}

func TestCreatePackage(t *testing.T) {
	p := dummyPackage()

	testCases := []struct {
		name       string
		body       PackageArgs
		buildStubs func(store *mocksdb.MockTxSubscriptionsStore)
		comparator func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "TESTOK",
			body: PackageArgs{
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			},
			buildStubs: func(store *mocksdb.MockTxSubscriptionsStore) {
				store.EXPECT().CreatePackage(gomock.Any(), db.CreatePackageParams{
					Name:        p.Name,
					Description: p.Description,
					Price:       p.Price,
				}).Times(1).Return(p, nil)
			},
			comparator: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				comparePostRequest(t, recorder.Body, p)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctlr := gomock.NewController(t)

			store := mocksdb.NewMockTxSubscriptionsStore(ctlr)
			testCase.buildStubs(store)
			defer ctlr.Finish()

			jsonPayload, err := json.Marshal(testCase.body)

			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/subscriptions", bytes.NewReader(jsonPayload))

			require.NoError(t, err)

			res := httptest.NewRecorder()

			srv := newTestServer(store)

			srv.router.ServeHTTP(res, req)

			testCase.comparator(t, res)

		})
	}
}

type GetPackagesResponse struct {
	Packages []db.Package `json:"packages"`
}

func compareResponses(t *testing.T, body *bytes.Buffer, packages []db.Package) {
	var response GetPackagesResponse

	jsonPayload, err := io.ReadAll(body)

	require.NoError(t, err)

	err = json.Unmarshal(jsonPayload, &response)

	require.NoError(t, err)
	require.Equal(t, response.Packages, packages)
}

type GetPackageResponse struct {
	Package db.Package `json:"package"`
}

func comparePostRequest(t *testing.T, body *bytes.Buffer, p db.Package) {
	var packageResponse GetPackageResponse
	data, err := io.ReadAll(body)

	require.NoError(t, err)

	err = json.Unmarshal(data, &packageResponse)

	require.NoError(t, err)
	require.Equal(t, packageResponse.Package, p)
}
