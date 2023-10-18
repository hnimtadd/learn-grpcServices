package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"grpcCource/client"
	"grpcCource/pb"
	"grpcCource/sample"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func testCreateLaptop(laptopClient *client.LaptopClient) {
	laptopClient.CreateLaptop(sample.NewLaptop())
}

func testSearchLaptop(laptopClient *client.LaptopClient) {
	for i := 0; i < 10; i++ {
		laptop := sample.NewLaptop()
		laptopClient.CreateLaptop(laptop)
	}

	filter := &pb.Filter{
		MaxPriceUsd: 3000,
		MinCpuCores: 4,
		MinCpuGhz:   2.0,
		MinRam: &pb.Memory{
			Unit:  pb.Memory_GIGABYTE,
			Value: 10,
		},
	}
	laptopClient.SearchLaptop(filter)
}

func testUploadImage(laptopClient *client.LaptopClient) {
	laptop := sample.NewLaptop()
	laptopClient.CreateLaptop(laptop)
	imagePath := "./sample/image.jpg"
	laptopClient.UploadImage(laptop.GetId(), imagePath)
}

func testRating(laptopClient *client.LaptopClient) {
	laptopIDs := []string{}
	scores := []float64{}
	for i := 0; i < 10; i++ {
		laptop := sample.NewLaptop()
		laptopClient.CreateLaptop(laptop)
		laptopIDs = append(laptopIDs, laptop.Id)
		scores = append(scores, sample.NewScore())
	}
	laptopClient.RateLaptop(laptopIDs, scores)
}

const (
	username        = "admin"
	password        = "secret"
	refreshDuration = 3 * time.Second
)

func authMethods() map[string]bool {
	const laptopServicePath = "/grpcCourse.pcbook.LaptopService/"
	return map[string]bool{
		laptopServicePath + "CreateLaptop": true,
		laptopServicePath + "UploadImage":  true,
		laptopServicePath + "RateLaptop":   true,
	}
}

func loadTLSCredenditial() (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair("./cert/client-cert.pem", "./cert/client-key.pem")
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
		RootCAs:      certPool,
	}
	return credentials.NewTLS(config), nil
}

func main() {
	serverAddress := flag.String("address", "", "the server address")
	enableTLS := flag.Bool("tls", false, "enable SSL/TLS")
	flag.Parse()

	log.Printf("dial server %s, TLS = %t", *serverAddress, *enableTLS)

	grpcOptions := []grpc.DialOption{}
	if *enableTLS {
		cred, err := loadTLSCredenditial()
		if err != nil {
			log.Fatal("Cannot load tls credential of CA", err.Error())
		}
		grpcOptions = append(grpcOptions, grpc.WithTransportCredentials(cred))
	} else {
		grpcOptions = append(grpcOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	cc1, err := grpc.Dial(*serverAddress, grpcOptions...)
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	authClient := client.NewAuthClient(cc1, username, password)
	interceptor, err := client.NewAuthInterceptor(authClient, authMethods(), refreshDuration)
	if err != nil {
		log.Fatal("Cannot create auth interceptor: ", err.Error())
	}
	grpcOptions = append(
		grpcOptions,
		grpc.WithChainUnaryInterceptor(
			interceptor.Unary(),
		),
		grpc.WithChainStreamInterceptor(
			interceptor.Stream(),
		),
	)
	cc2, err := grpc.Dial(*serverAddress, grpcOptions...)
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}
	laptopClient := client.NewLaptopClient(cc2, username, password)

	testCreateLaptop(laptopClient)
	// testSearchLaptop(laptopClient)
	// testUploadImage(laptopClient)
	// testRating(laptopClient)
}
