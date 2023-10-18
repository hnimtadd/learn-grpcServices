package client

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthInterceptor struct {
	authClient  *AuthClient
	authMethods map[string]bool
	accessToken string
}

func NewAuthInterceptor(
	authClient *AuthClient,
	authMethods map[string]bool,
	refreshDuration time.Duration,
) (*AuthInterceptor, error) {
	intercepter := &AuthInterceptor{
		authClient:  authClient,
		authMethods: authMethods,
	}
	err := intercepter.scheduleRefreshToken(refreshDuration)
	if err != nil {
		return nil, err
	}
	return intercepter, nil
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req,
		reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		log.Printf("--> unary interceptor: %s", method)
		if interceptor.authMethods[method] {
			log.Printf("with credentials")
			return invoker(interceptor.attachToken(ctx), method, req, reply, cc, opts...)
		}
		log.Printf("with no credentials")
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (interceptor *AuthInterceptor) Stream() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		log.Printf("--> stream interceptor: %s", method)
		if interceptor.authMethods[method] {
			log.Printf("with credentials")
			return streamer(interceptor.attachToken(ctx), desc, cc, method, opts...)
		}

		log.Printf("with no credentials")
		return streamer(ctx, desc, cc, method, opts...)
	}
}
func (interceptor *AuthInterceptor) attachToken(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", interceptor.accessToken)
}

func (intercepter *AuthInterceptor) refreshToken() error {
	token, err := intercepter.authClient.Login()
	if err != nil {
		return err
	}
	log.Printf("Token refreshed: %v", token)
	intercepter.accessToken = token
	return nil
}

func (intercepter *AuthInterceptor) scheduleRefreshToken(refreshDuration time.Duration) error {
	err := intercepter.refreshToken()
	if err != nil {
		return err
	}
	go func() {
		wait := time.NewTicker(refreshDuration)
		for {
			select {
			case <-wait.C:
				err := intercepter.refreshToken()
				if err != nil {
					wait.Stop()
					time.Sleep(1)
					wait.Reset(refreshDuration)
				}
			}
		}
	}()
	return nil
}
