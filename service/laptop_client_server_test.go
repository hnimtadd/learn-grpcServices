package service_test

import (
	"context"
	"grpcCource/pkg/pb"
	"grpcCource/pkg/serializer"
	"grpcCource/sample"
	"grpcCource/service"
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestClientCreateLaptop(t *testing.T) {
	t.Parallel()
	laptopStore := service.NewInMemoryLaptopStore()
	laptopServer, serverAddress := startTestLaptopServer(t, laptopStore, nil, nil)
	laptopClient := newTestLaptopClient(t, serverAddress)
	laptop := sample.NewLaptop()
	expectedID := laptop.Id
	rsp, err := laptopClient.CreateLaptop(context.TODO(), &pb.CreateLaptopRequest{
		Laptop: laptop,
	})
	require.NoError(t, err)
	require.NotNil(t, rsp)
	require.Equal(t, rsp.Id, expectedID)

	// check that the laptopReally save in server
	other, err := laptopServer.LaptopStore.Find(laptop.Id)
	require.NoError(t, err)
	require.NotNil(t, other)
	requireSameLaptop(t, laptop, other)
}

func TestClientSearchLaptop(t *testing.T) {
	t.Parallel()
	laptopStore := service.NewInMemoryLaptopStore()
	_, serverAddress := startTestLaptopServer(t, laptopStore, nil, nil)
	laptopClient := newTestLaptopClient(t, serverAddress)
	filter := &pb.Filter{
		MaxPriceUsd: 3000,
		MinCpuCores: 2,
		MinCpuGhz:   2.2,
		MinRam:      &pb.Memory{Value: 8, Unit: pb.Memory_GIGABYTE},
	}
	expectedID := map[string]bool{}

	for i := 0; i < 6; i++ {
		laptop := sample.NewLaptop()
		switch i {
		case 0:
			laptop.PriceUsd = 5000
		case 1:
			laptop.Cpu.NumberCores = 1
		case 2:
			laptop.Cpu.MinGhz = 2.0
		case 3:
			laptop.Ram = &pb.Memory{Value: 8, Unit: pb.Memory_BYTE}
		case 4:
			laptop.PriceUsd = 2500
			laptop.Cpu.NumberCores = 4
			laptop.Cpu.MinGhz = 2.8
			laptop.Ram = &pb.Memory{Value: 9, Unit: pb.Memory_GIGABYTE}
			expectedID[laptop.Id] = true
		case 5:
			laptop.PriceUsd = 2600
			laptop.Cpu.NumberCores = 4
			laptop.Cpu.MinGhz = 2.9
			laptop.Ram = &pb.Memory{Value: 16, Unit: pb.Memory_GIGABYTE}
			expectedID[laptop.Id] = true
		}
		rsp, err := laptopClient.CreateLaptop(context.Background(), &pb.CreateLaptopRequest{Laptop: laptop})
		require.NoError(t, err)
		require.NotNil(t, rsp)
		require.Equal(t, rsp.GetId(), laptop.GetId())
	}

	stream, err := laptopClient.SearchLaptop(context.Background(), &pb.SearchLaptopRequest{Filter: filter})
	require.NoError(t, err)
	require.NotNil(t, stream)
	found := 0

	for {
		rsp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		require.NotNil(t, rsp)

		require.Contains(t, expectedID, rsp.GetLaptop().GetId())
		found += 1
	}
	require.Equal(t, found, len(expectedID))
}

func startTestLaptopServer(t *testing.T, laptopStore service.LaptopStore, imageStore service.ImageStore, ratingStore service.RatingStore) (*service.LaptopServer, string) {
	laptopServer := service.NewLaptopServer(laptopStore, imageStore, ratingStore)
	grpcServer := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	go grpcServer.Serve(listener)
	return laptopServer, listener.Addr().String()
}

func newTestLaptopClient(t *testing.T, serverAddress string) pb.LaptopServiceClient {
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	return pb.NewLaptopServiceClient(conn)
}

func requireSameLaptop(t *testing.T, laptop1 *pb.Laptop, laptop2 *pb.Laptop) {
	json1, err := serializer.ProtoBufToJSON(laptop1)
	require.NoError(t, err)
	json2, err := serializer.ProtoBufToJSON(laptop2)
	require.NoError(t, err)
	require.Equal(t, json1, json2)
}
