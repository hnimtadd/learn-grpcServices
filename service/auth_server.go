package service

import (
	"context"
	"grpcCource/pkg/pb"
	"grpcCource/pkg/store"
	"grpcCource/pkg/token"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	userStore store.UserStore
	jwtManage *token.JWTManager
	pb.UnimplementedAuthServiceServer
}

func NewAuthServer(userStore store.UserStore, jwtManage *token.JWTManager) *AuthServer {
	authServer := &AuthServer{
		userStore: userStore,
		jwtManage: jwtManage,
	}
	return authServer
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var (
		userName = req.GetUserName()
		password = req.GetPassword()
	)
	user, err := s.userStore.Find(userName)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "User with username %v not found", userName)
	}
	if user == nil || !user.VerifyPassword(password) {
		return nil, status.Error(codes.NotFound, "Username or password incorrect")
		// Usernot loging
	}
	token, err := s.jwtManage.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cannot generate token for username: %v", err)
	}
	rsp := &pb.LoginResponse{Token: token}
	return rsp, nil
}
