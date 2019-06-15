package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net"

	maximum "github.com/ReigoHein/udemy-grpc-course/go-grpc-find-maximum/models"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) FindMaximum(stream maximum.FindMaximumService_FindMaximumServer) error {
	fmt.Printf("FindMaximum function was invoked with a streaming request\n")
	result := int32(math.MinInt32)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}

		number := req.GetNumber().GetNumber()
		if number > result {
			result = number

			sendErr := stream.Send(&maximum.MaximumResponse{
				Result: result,
			})
			if sendErr != nil {
				log.Fatalf("Error while sending data to client: %v", sendErr)
				return sendErr
			}
		}
	}
}

func main() {
	fmt.Println("Server started")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)
	maximum.RegisterFindMaximumServiceServer(s, &server{})

	fmt.Println("Starting to serve on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
