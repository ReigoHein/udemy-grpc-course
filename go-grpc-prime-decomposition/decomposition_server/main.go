package main

import (
	"fmt"
	"log"
	"net"
	"time"

	primedecomp "github.com/ReigoHein/udemy-grpc-course/go-grpc-prime-decomposition/models"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Decompose(req *primedecomp.DecompositionRequest, stream primedecomp.PrimeDecompositionService_DecomposeServer) error {
	fmt.Printf("Prime decompose was invoked with: %+v\n", req)

	number := req.Decomposition.GetNumber()
	divider := 2

	for number > 1 {
		remainder := number % int32(divider)

		if remainder == 0 {
			result := &primedecomp.DecompositionResponse{
				Multiplier: int32(divider),
			}
			stream.Send(result)

			time.Sleep(200 * time.Millisecond)
			number = number / int32(divider)
		} else {
			divider = divider + 1
		}
	}
	return nil
}

func main() {
	fmt.Println("Server started")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)
	primedecomp.RegisterPrimeDecompositionServiceServer(s, &server{})

	fmt.Println("Starting to serve on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
