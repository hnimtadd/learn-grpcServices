package service

import (
	"context"
	"grpcCource/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	userStore UserStore
	jwtManage *JWTManager
	pb.UnimplementedAuthServiceServer
}

func NewAuthServer(userStore UserStore, jwtManage *JWTManager) *AuthServer {
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
