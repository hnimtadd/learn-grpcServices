package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"grpcCource/pkg/models"
	"grpcCource/pkg/pb"
	"grpcCource/pkg/store"
	"grpcCource/pkg/token"
	"grpcCource/service"
	"grpcCource/utils"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func unaryInterceptorHandler() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		log.Printf("--> unary: %v", info.FullMethod)
		return handler(ctx, req)
	}
}

func streamInterceptorHandler() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		log.Printf("--> stream: %v", info.FullMethod)
		return handler(srv, ss)
	}
}

const (
	secretKey     = "mysecret"
	tokenDuration = time.Minute * 5
)

func accessibleRoles() map[string][]string {
	const laptopServicePath = "/grpcCourse.pcbook.LaptopService/"

	return map[string][]string{
		laptopServicePath + "CreateLaptop": {"admin"},
		laptopServicePath + "UploadLaptop": {"admin"},
		laptopServicePath + "SearchLaptop": {"admin", "user"},
		laptopServicePath + "RateLaptop":   {"admin", "user"},
	}
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair("./cert/server-cert.pem", "./cert/server-key.pem")
	if err != nil {
		return nil, err
	}

	pemServerCA, err := os.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("Failed to add server CA's certificate")
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}
	return credentials.NewTLS(config), nil
}

func createUser(userStore store.UserStore, username, password, role string) error {
	user, err := models.NewUser(username, password, role)
	if err != nil {
		return err
	}
	return userStore.Add(user)
}

func seedUsers(userStore store.UserStore) error {
	if err := createUser(userStore, "admin", "secret", "admin"); err != nil {
		return err
	}
	return createUser(userStore, "user1", "secret", "user")
}

func main() {
	fmt.Println("Hello world from server")
	port := flag.Int("port", 0, "the server port")
	enableTLS := flag.Bool("tls", false, "enable SSL/TLS")
	dsn := flag.String("dsn", "", "dsn of mongo")
	database := flag.String("db", "", "database of mongo")
	flag.Parse()
	log.Printf("Server listen on port %d\n, TLS=%t", *port, *enableTLS)

	var (
		laptopStore = service.NewInMemoryLaptopStore()
		imageStore  = service.NewDickImageStore("images")
		ratingStore = service.NewInMemoryRatingStore()
	)
	if *database != "" && *dsn != "" {
		db, err := utils.NewMongoDB(*dsn, *database)
		if err != nil {
			log.Fatal("Cannot start mongo connection:", err)
		}
		ratingStore = service.NewMongoScoreStore(db)
	}

	laptopServer := service.NewLaptopServer(
		laptopStore,
		imageStore,
		ratingStore,
	)
	userStore := service.NewInMemoryUserStore()
	if err := seedUsers(userStore); err != nil {
		log.Fatalf("Cannot seed user: %v", err)
	}

	jwtManager := token.NewJWTManager(secretKey, tokenDuration)
	interceptor := service.NewAuthInterceptor(jwtManager, accessibleRoles())
	grpcOptions := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptor.Unary(),
		),
		grpc.ChainStreamInterceptor(
			interceptor.Stream(),
		),
	}

	if *enableTLS {
		cred, err := loadTLSCredentials()
		if err != nil {
			log.Fatal("Cannot load tls credents", err)
		}
		grpcOptions = append(grpcOptions, grpc.Creds(cred))
	}

	authServer := service.NewAuthServer(userStore, jwtManager)
	grpcServer := grpc.NewServer(grpcOptions...)
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	pb.RegisterAuthServiceServer(grpcServer, authServer)
	reflection.Register(grpcServer)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Cannot start server: ", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Cannot start server: ", err)
	}
}
