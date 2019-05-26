package main

import (
	"context"
	"log"
	"net"

	sum "github.com/ReigoHein/udemy-grpc-course/go-grpc-sum/models"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) Add(ctx context.Context, in *sum.AddRequest) (*sum.AddResponse, error) {
	log.Printf("Received %+v", in)
	numbers := in.GetNumbers()

	result := int32(0)

	for _, nr := range numbers {
		result += nr
	}

	return &sum.AddResponse{Result: result}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	sum.RegisterSumServer(s, &server{})

	log.Printf("starting serving on: %v", port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
