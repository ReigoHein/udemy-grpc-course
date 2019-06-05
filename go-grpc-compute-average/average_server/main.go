package main

import (
	"fmt"
	"io"
	"log"
	"net"

	average "github.com/ReigoHein/udemy-grpc-course/go-grpc-compute-average/models"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) ComputeAverage(stream average.ComputeAverageService_ComputeAverageServer) error {
	fmt.Printf("Compute average was invoked with a streaming request")

	sum := int32(0)
	total := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&average.AverageResponse{
				Result: float32(sum) / float32(total),
			})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		number := req.GetAverage().GetNumber()

		sum += number
		total++
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

	average.RegisterComputeAverageServiceServer(s, &server{})

	fmt.Println("Starting to serve on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
