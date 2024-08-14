package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mockdb "github.com/CodeSingerGnC/MicroBank/db/mock"
	db "github.com/CodeSingerGnC/MicroBank/db/sqlc"
	"github.com/CodeSingerGnC/MicroBank/util"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type MockResult struct{
	ID 		int64
	Rows	int64
}

func (r MockResult) LastInsertId() (int64, error) {
    return r.ID, nil
}

func (r MockResult) RowsAffected() (int64, error) {
    return r.ID, nil
}

type eqCreateUserParamsMatcher struct {
	arg db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPasswordHash(e.password, arg.HashPassword)
	if err != nil {
		return false
	}

	e.arg.HashPassword = arg.HashPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v ans password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher { 
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUser(t *testing.T) {
	user, password := randomUser(t)

	mockFailResult := MockResult{
		ID: 0,
		Rows: 0,
	}

	mockSuccessResult := MockResult{
		ID: 1,
		Rows: 1,
	}

	testCase := []struct {
		name			string
		body			gin.H
		buildStubs  	func(store *mockdb.MockStore)
		checkResponse	func(recoder *httptest.ResponseRecorder)
	} {
		{
			name: "OK",
			body: gin.H{
				"user_account": user.UserAccount,
				"password": password,
				"username": user.Username,
				"email": user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					UserAccount: user.UserAccount,
					Username: user.Username,
					Email: user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(mockSuccessResult, nil)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
				requireBodyMatchUser(t, recoder.Body, user)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"user_account": user.UserAccount,
				"password": user.HashPassword,
				"username": user.Username,
				"email": user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(mockFailResult, sql.ErrConnDone)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
		{
			name: "DuplicateUsername",
			body: gin.H{
				"user_account": user.UserAccount,
				"password": user.HashPassword,
				"username": user.Username,
				"email": user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(mockFailResult, &mysql.MySQLError{Number: 1062})
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recoder.Code)
			},
		},
		{
			name: "InvalidUserAccount",
			body: gin.H{
				"user_account": "invalid-user-acount#1",
				"password": user.HashPassword,
				"username": user.Username,
				"email": user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"user_account": user.UserAccount,
				"password": user.HashPassword,
				"username": user.Username,
				"email": "invalid-email",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"user_account": user.UserAccount,
				"password": "123",
				"username": user.Username,
				"email": user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
	}

	for i := range testCase {
		tc := testCase[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recoder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recoder, request)
			tc.checkResponse(recoder)
		})
	}
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		UserAccount: util.RandomUser(),
		HashPassword: hashPassword,
		Username: util.RandomString(6),
		Email: util.RandomEmail(),
	}
	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	var userAccount string
	err = json.Unmarshal(data, &userAccount)
	require.NoError(t, err)

	require.Equal(t, userAccount, user.UserAccount)
}