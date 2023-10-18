package service

import (
	"bytes"
	"context"
	"errors"
	"grpcCource/pb"
	"io"
	"log"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LaptopServer struct {
	LaptopStore LaptopStore
	ImageStore  ImageStore
	RatingStore RatingStore
	pb.UnimplementedLaptopServiceServer
}

func NewLaptopServer(laptopStore LaptopStore, imageStore ImageStore, RatingStore RatingStore) *LaptopServer {
	return &LaptopServer{
		LaptopStore: laptopStore,
		ImageStore:  imageStore,
		RatingStore: RatingStore,
	}
}

func (s *LaptopServer) CreateLaptop(
	ctx context.Context,
	req *pb.CreateLaptopRequest,
) (*pb.CreateLaptopResponse, error) {
	laptop := req.GetLaptop()
	log.Printf("received a create-laptop request with id: %v", laptop.Id)
	if len(laptop.Id) > 0 {
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "laptop ID is not a valid UUID: %v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new laptop ID: %v", err)
		}
		laptop.Id = id.String()
	}

	if err := contextError(ctx); err != nil {
		return nil, status.Errorf(codes.Unavailable, "Service ended: %v", err.Error())
	}
	// save the laptop to store
	err := s.LaptopStore.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "Cannot save laptop to the store: %v", err)
	}
	log.Printf("Saved laptop with id: %v\n", laptop.Id)
	return &pb.CreateLaptopResponse{
		Id: laptop.Id,
	}, nil
}

func (s *LaptopServer) SearchLaptop(
	req *pb.SearchLaptopRequest,
	stream pb.LaptopService_SearchLaptopServer,
) error {
	var (
		filter   = req.GetFilter()
		callback = func(laptop *pb.Laptop) error {
			msg := &pb.SearchLaptopResponse{
				Laptop: laptop,
			}
			err := stream.SendMsg(msg)
			if err != nil {
				log.Printf("Cannot send message to server: %v", err)
				return err
			}
			log.Printf("Sended laptop with id: %v", laptop.GetId())
			return nil
		}
	)

	err := s.LaptopStore.Search(stream.Context(), filter, callback)
	if err != nil {
		log.Printf("Cannot search laptop: %v", err)
		return status.Errorf(codes.Internal, "Unexpected error: %v", err)
	}
	return nil
}

func (s *LaptopServer) UploadImage(stream pb.LaptopService_UploadImageServer) error {
	msg, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Cann't received upload image info: %v", err)
	}
	imgInfo := msg.GetInfo()
	if imgInfo == nil {
		return status.Error(codes.InvalidArgument, "First message should be info of image")
	}
	var (
		laptopId  = imgInfo.GetLaptopId()
		imageType = imgInfo.GetImageType()
		dataChunk = bytes.Buffer{}
		lenReaded = 0
	)

	laptop, err := s.LaptopStore.Find(laptopId)
	if err != nil {
		return status.Errorf(codes.Internal, "Cannot find laptop with id: %v", laptopId)
	}
	if laptop == nil {
		return status.Errorf(codes.NotFound, "Laptop not found: %v", laptopId)
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// Complete client stream, save
				break
			}
			return status.Errorf(codes.InvalidArgument, "Cann't received upload image info: %v", err)
		}
		chunk := msg.GetChunkData()
		if chunk == nil {
			return status.Error(codes.InvalidArgument, "Data chunk not found")
		}
		r, err := dataChunk.Write(chunk)
		if err != nil {
			return status.Errorf(codes.Internal, "Cannot read byte from chunk: %v", err)
		}
		lenReaded += r
	}

	imageID, err := s.ImageStore.Save(laptopId, imageType, dataChunk)
	if err != nil {
		return status.Errorf(codes.Internal, "Cannot save image to store: %v", err)
	}
	rsp := &pb.UploadImageResponse{
		Id:   imageID,
		Size: uint32(lenReaded),
	}
	if err := stream.SendAndClose(rsp); err != nil {
		return status.Errorf(codes.Internal, "Unexpected error: %v", err)
	}
	return nil
}

func (s *LaptopServer) RateLaptop(stream pb.LaptopService_RateLaptopServer) error {
	reqCh := make(chan *pb.RateLaptopRequest, 10)
	rspCh := make(chan *pb.RateLaptopResponse, 10)
	errorCh := make(chan error, 1)
	// Receving side
	go func() {
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				log.Println("No more data from client")
				return
			}
			if err != nil {
				errorCh <- status.Errorf(codes.Internal, "Cannot read message from client: %v", err)
				return
			}
			time.Sleep(time.Second)
			log.Printf("Received request from client: %v", req.String())
			reqCh <- req
		}
	}()

	// Sending side
	go func() {
		for {
			select {
			case req := <-reqCh:
				var (
					laptopID = req.GetLaptopId()
					score    = req.GetScore()
				)
				laptop, err := s.LaptopStore.Find(laptopID)
				if err != nil {
					errorCh <- err
					return
				}
				if laptop == nil {
					errorCh <- status.Errorf(codes.NotFound, "Cannot found laptop with id: %v", laptopID)
					return
				}
				if score < 0 || score > 10 {
					errorCh <- status.Errorf(codes.InvalidArgument, "Laptop score must be from 0 - 10, your score: %v", score)
					return
				}
				rating, err := s.RatingStore.Add(laptopID, score)
				if err != nil {
					errorCh <- err
					return
				}
				if rating == nil {
					errorCh <- status.Errorf(codes.NotFound, "Cannot found rating of laptop with id: %v", laptopID)
					return
				}
				rsp := &pb.RateLaptopResponse{
					LaptopId:     laptopID,
					RatedCount:   rating.Count,
					AverageScore: float64((rating.Count) / rating.Count),
				}
				rspCh <- rsp
			}
		}
	}()

	for {
		select {
		case <-stream.Context().Done():
			log.Printf("context done: %v", stream.Context().Err())
			return stream.Context().Err()
		case err := <-errorCh:
			return err
		case rsp := <-rspCh:
			log.Printf("New response to client")
			if err := stream.Send(rsp); err != nil {
				log.Printf("Cannot send response to server: %v", err)
				return err
			}
		}
	}
}

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return ctx.Err()
	case context.DeadlineExceeded:
		return ctx.Err()
	default:
		return nil
	}
}
