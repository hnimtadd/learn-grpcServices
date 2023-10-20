package service_test

import (
	"context"
	"grpcCource/pkg/models"
	"grpcCource/pkg/pb"
	"grpcCource/pkg/store"
	"grpcCource/pkg/token"
	"grpcCource/service"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuthServer(t *testing.T) {
	t.Parallel()

	var (
		username        = "hnimtadd"
		password        = "hnimtadd"
		role            = "user"
		secretKey       = "secret"
		refreshDuration = time.Second * 5
	)
	user, err := models.NewUser(username, password, role)
	require.NoError(t, err)
	require.NotNil(t, user)

	userStore := service.NewInMemoryUserStore()
	userStore.Add(user)
	require.NotNil(t, userStore)
	jwtManager := token.NewJWTManager(secretKey, refreshDuration)
	testCases := []struct {
		name       string
		userStore  store.UserStore
		username   string
		password   string
		jwtManager *token.JWTManager
		code       codes.Code
	}{
		{
			name:       "login success",
			userStore:  userStore,
			username:   username,
			password:   password,
			code:       codes.OK,
			jwtManager: jwtManager,
		},
		{
			name:       "user not found",
			username:   "username",
			password:   password,
			userStore:  userStore,
			jwtManager: jwtManager,
			code:       codes.NotFound,
		}, {
			name:       "Password failed",
			username:   username,
			password:   "password",
			userStore:  userStore,
			jwtManager: jwtManager,
			code:       codes.NotFound,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := &pb.LoginRequest{
				UserName: tc.username,
				Password: tc.password,
			}
			server := service.NewAuthServer(tc.userStore, tc.jwtManager)
			rsp, err := server.Login(context.TODO(), req)
			if tc.code == codes.OK {
				require.NoError(t, err)
				require.NotNil(t, rsp)
				require.NotEmpty(t, rsp.Token)
			} else {
				require.Error(t, err)
				require.Nil(t, rsp)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tc.code, st.Code())
			}
		},
		)
	}
}
