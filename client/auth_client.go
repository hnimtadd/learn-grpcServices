package client

import (
	"context"
	"grpcCource/pb"
	"time"

	"google.golang.org/grpc"
)

type AuthClient struct {
	service  pb.AuthServiceClient
	userName string
	password string
}

func NewAuthClient(cc *grpc.ClientConn, username string, password string) *AuthClient {
	service := pb.NewAuthServiceClient(cc)
	return &AuthClient{
		service:  service,
		userName: username,
		password: password,
	}
}

func (client *AuthClient) Login() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &pb.LoginRequest{UserName: client.userName, Password: client.password}
	rsp, err := client.service.Login(ctx, req)
	if err != nil {
		return "", err
	}
	return rsp.GetToken(), nil
}
