package service

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	jwtManager      *JWTManager
	accessibleRoles map[string][]string
}

func NewAuthInterceptor(jwtManager *JWTManager, accessibleRoles map[string][]string) *AuthInterceptor {
	return &AuthInterceptor{
		jwtManager:      jwtManager,
		accessibleRoles: accessibleRoles,
	}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		log.Printf("-->unary, auth: %v", info.FullMethod)
		if err := interceptor.authorize(ctx, info.FullMethod); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func (interceptor *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		log.Printf("-->stream, auth: %v", info.FullMethod)
		if err := interceptor.authorize(ss.Context(), info.FullMethod); err != nil {
			return err
		}
		return handler(srv, ss)
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, grpcMethod string) error {
	roles, ok := interceptor.accessibleRoles[grpcMethod]
	if !ok {
		// default will be public path
		log.Println("with no credential")
		return nil
	}
	log.Println("with credential")

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "metadata is not provided")
	}
	values, ok := md["authorization"]
	if !ok || len(values) == 0 {
		return status.Error(codes.Unauthenticated, "authorization token is not provide")
	}
	token := values[0]
	claims, err := interceptor.jwtManager.Verify(token)
	if err != nil {
		return status.Error(codes.Unauthenticated, "cannot verify token")
	}
	for _, role := range roles {
		if claims.Role == role {
			log.Println("authorized")
			return nil
		}
	}
	return status.Error(codes.PermissionDenied, "no permission to access this RPC")
}
