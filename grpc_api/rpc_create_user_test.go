package grpc_api

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/mnakhaev/simplebank/db/mock"
	db "github.com/mnakhaev/simplebank/db/sqlc"
	"github.com/mnakhaev/simplebank/pb"
	"github.com/mnakhaev/simplebank/util"
	"github.com/mnakhaev/simplebank/worker"
	mockwk "github.com/mnakhaev/simplebank/worker/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type eqCreateUserTxParamsMatcher struct {
	arg      db.CreateUserTxParams
	password string
	user     db.User
}

func (expected eqCreateUserTxParamsMatcher) Matches(x interface{}) bool {
	actualArg, ok := x.(db.CreateUserTxParams)
	if !ok {
		return false
	}
	if err := util.CheckPassword(expected.password, actualArg.HashedPassword); err != nil {
		return false
	}

	// some hack
	expected.arg.HashedPassword = actualArg.HashedPassword

	if !reflect.DeepEqual(expected.arg.CreateUserParams, actualArg.CreateUserParams) {
		return false
	}

	if err := actualArg.AfterCreate(expected.user); err != nil {
		return false
	}
	return true
}

func (expected eqCreateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", expected.arg, expected.password)
}

func EqCreateUserTxParams(arg db.CreateUserTxParams, password string, user db.User) gomock.Matcher {
	return eqCreateUserTxParamsMatcher{arg, password, user}
}
func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)
	testCases := []struct {
		name          string
		req           *pb.CreateUserRequest
		buildStubs    func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor)
		checkResponse func(t *testing.T, res *pb.CreateUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor) {
				arg := db.CreateUserTxParams{
					CreateUserParams: db.CreateUserParams{
						Username: user.Username,
						FullName: user.FullName,
						Email:    user.Email,
					},
				}
				store.EXPECT().
					CreateUserTx(gomock.Any(), EqCreateUserTxParams(arg, password, user)).
					Times(1).
					Return(db.CreateUserTxResult{User: user}, nil)

				taskPayload := &worker.PayloadSendVerifyEmail{Username: user.Username}
				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), taskPayload, gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				createdUser := res.GetUser()
				require.Equal(t, user.Username, createdUser.Username)
				require.Equal(t, user.FullName, createdUser.FullName)
				require.Equal(t, user.Email, createdUser.Email)
			},
		},
		{
			name: "InternalError",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},

			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor) {
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateUserTxResult{}, sql.ErrConnDone)
				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				// TODO: check
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
		// {
		// 	name: "Duplicate username",
		// 	req: &pb.CreateUserRequest{
		// 		Username: user.Username,
		// 		Password: password,
		// 		FullName: user.FullName,
		// 		Email:    user.Email,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		arg := db.CreateUserTxParams{
		// 			CreateUserParams: db.CreateUserParams{
		// 				Username: user.Username,
		// 				FullName: user.FullName,
		// 				Email:    user.Email,
		// 			},
		// 		}
		// 		store.EXPECT().CreateUserTx(gomock.Any(), EqCreateUserTxParams(arg, password)).Times(1).Return(db.User{}, &pq.Error{Code: "12323"})
		// 	},
		// 	checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
		// 		// TODO: check
		// 		require.NoError(t, err)
		// 		require.NotNil(t, res)
		// 		createdUser := res.GetUser()
		// 		require.Equal(t, user.Username, createdUser.Username)
		// 		require.Equal(t, user.FullName, createdUser.FullName)
		// 		require.Equal(t, user.Email, createdUser.Email)
		// 	},
		// },
		// {
		// 	name: "invalid username",
		// 	req: &pb.CreateUserRequest{
		// 		Username: "invalid-user1",
		// 		Password: password,
		// 		FullName: user.FullName,
		// 		Email:    user.Email,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().CreateUserTx(gomock.Any(), gomock.Any()).Times(0)
		// 	},
		// 	checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
		// 		// TODO: check
		// 		require.NoError(t, err)
		// 		require.NotNil(t, res)
		// 		createdUser := res.GetUser()
		// 		require.Equal(t, user.Username, createdUser.Username)
		// 		require.Equal(t, user.FullName, createdUser.FullName)
		// 		require.Equal(t, user.Email, createdUser.Email)
		// 	},
		// },
		// {
		// 	name: "invalid email",
		// 	req: &pb.CreateUserRequest{
		// 		Username: user.Username,
		// 		Password: password,
		// 		FullName: user.FullName,
		// 		Email:    "123",
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().CreateUserTx(gomock.Any(), gomock.Any()).Times(0)
		// 	},
		// 	checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
		// 		// TODO: check
		// 		require.NoError(t, err)
		// 		require.NotNil(t, res)
		// 		createdUser := res.GetUser()
		// 		require.Equal(t, user.Username, createdUser.Username)
		// 		require.Equal(t, user.FullName, createdUser.FullName)
		// 		require.Equal(t, user.Email, createdUser.Email)
		// 	},
		// },
		// {
		// 	name: "too short password",
		// 	req: &pb.CreateUserRequest{
		// 		Username: user.Username,
		// 		Password: "123",
		// 		FullName: user.FullName,
		// 		Email:    user.Email,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().CreateUserTx(gomock.Any(), gomock.Any()).Times(0)
		// 	},
		// 	checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
		// 		// TODO: check
		// 		require.NoError(t, err)
		// 		require.NotNil(t, res)
		// 		createdUser := res.GetUser()
		// 		require.Equal(t, user.Username, createdUser.Username)
		// 		require.Equal(t, user.FullName, createdUser.FullName)
		// 		require.Equal(t, user.Email, createdUser.Email)
		// 	},
		// },
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()

			taskCtrl := gomock.NewController(t)
			defer taskCtrl.Finish()

			store := mockdb.NewMockStore(storeCtrl)
			taskDistributor := mockwk.NewMockTaskDistributor(taskCtrl)
			tc.buildStubs(store, taskDistributor)

			server := newTestServer(t, store, taskDistributor)
			res, err := server.CreateUser(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func randomUser(t *testing.T) (db.User, string) {
	psw := util.RandomString(10)
	hp, err := util.HashPassword(psw)
	require.NoError(t, err)
	return db.User{
		Username:       util.RandomString(10),
		HashedPassword: hp,
		FullName:       util.RandomString(10),
		Email:          util.RandomString(10) + "@email.com",
	}, psw
}
