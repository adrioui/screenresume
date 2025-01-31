package services

import (
	"context"
	"log"
	"screenresume/internal/repositories"
	pb "screenresume/internal/repositories"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type ScreenResumeService struct {
	client pb.ResumeScreenerClient
}

func NewScreenResumeService(conn *grpc.ClientConn) ScreenResumeService {
	client := repositories.NewResumeScreenerClient(conn)
	return ScreenResumeService{client: client}
}

func (s ScreenResumeService) ScreenResume(request *pb.ScreenResumeRequest) (*pb.ScreenResumeResponse, error) {
	response, err := s.client.ScreenResume(context.Background(), request)
	if err != nil {
		st, _ := status.FromError(err)
		log.Printf("error: %v", st)
		return nil, err
	}
	return response, nil
}
