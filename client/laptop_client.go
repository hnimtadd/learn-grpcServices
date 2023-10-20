package client

import (
	"bufio"
	"context"
	"grpcCource/pkg/pb"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LaptopClient struct {
	service pb.LaptopServiceClient
}

func NewLaptopClient(cc *grpc.ClientConn, username string, password string) *LaptopClient {
	service := pb.NewLaptopServiceClient(cc)
	return &LaptopClient{
		service: service,
	}
}

func (client *LaptopClient) CreateLaptop(laptop *pb.Laptop) error {
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	rsp, err := client.service.CreateLaptop(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			log.Printf("Cannot create laptop: %v", st.Err().Error())
			return st.Err()
		}

		log.Printf("Cannot create laptop: %v", err)
		return err
	}
	log.Printf("Created laptop with id: %v", rsp.GetId())
	return nil

}
func (client *LaptopClient) SearchLaptop(filter *pb.Filter) error {
	req := &pb.SearchLaptopRequest{Filter: filter}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	stream, err := client.service.SearchLaptop(ctx, req)
	if err != nil {
		return err
	}
	for {
		rsp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Println("Server stream close")
				return nil
			}
			status, ok := status.FromError(err)
			if ok {
				log.Fatalf("Status from server, code: %v, err: %v", status.Code(), status.Err())
			} else {
				log.Fatalf("Error while recv from server: %v", err)
			}
			continue
		}
		laptop := rsp.GetLaptop()
		log.Printf("Received laptop from server: %s", laptop.String())
	}
}
func (client *LaptopClient) UploadImage(laptopId string, imagePath string) error {
	file, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	stream, err := client.service.UploadImage(ctx)
	if err != nil {
		return err
	}

	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Info{
			Info: &pb.ImageInfo{
				LaptopId:  laptopId,
				ImageType: filepath.Ext(imagePath),
			},
		},
	}
	if err := stream.Send(req); err != nil {
		return err
	}

	r := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	for {
		n, err := r.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}
		if err := stream.Send(req); err != nil {
			return err
		}
	}

	rsp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Printf("Received response from server:\nlaptopID: %v\nSize: %v", rsp.GetId(), rsp.GetSize())
	return nil
}
func (client *LaptopClient) RateLaptop(laptopIDs []string, scores []float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	stream, err := client.service.RateLaptop(ctx)
	if err != nil {
		return err
	}
	defer stream.CloseSend()
	errCh := make(chan error, 1)
	go func() {
		for i, laptopID := range laptopIDs {
			req := &pb.RateLaptopRequest{
				LaptopId: laptopID,
				Score:    scores[i],
			}
			if err := stream.Send(req); err != nil {
				errCh <- err
				log.Printf("Cannot send request to server: %v", err)
				return
			}
		}
		if err := stream.CloseSend(); err != nil {
			errCh <- err
		}
	}()
	go func() {
		for {
			rsp, err := stream.Recv()
			if err == io.EOF {
				log.Println("Server stop streaming")
				errCh <- nil
				return
			}
			if err != nil {
				log.Printf("Cannot receive response from server:%v", err)
				errCh <- err
				return
			}
			log.Printf("Received new response from server: %v", rsp.String())
		}
	}()

	for {
		select {
		case <-stream.Context().Done():
			log.Printf("Context done: %v", stream.Context().Err())
			return stream.Context().Err()
		case <-ctx.Done():
			log.Printf("client Context done: %v", ctx.Err())
			return ctx.Err()
		case err := <-errCh:
			return err
		}
	}
}
