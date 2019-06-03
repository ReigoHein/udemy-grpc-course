package main

import (
	"context"
	"fmt"
	"io"
	"log"

	primedecomp "github.com/ReigoHein/udemy-grpc-course/go-grpc-prime-decomposition/models"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Started prime decomposition client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := primedecomp.NewPrimeDecompositionServiceClient(cc)

	decomposeNumber(c, 120)
	decomposeNumber(c, 917)
	decomposeNumber(c, 971)
}

func decomposeNumber(c primedecomp.PrimeDecompositionServiceClient, number int) {
	fmt.Printf("Decomposing prime decomposing: %v\n", number)

	req := &primedecomp.DecompositionRequest{
		Decomposition: &primedecomp.Decomposition{
			Number: int32(number),
		},
	}

	resStream, err := c.Decompose(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Decompose gRPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("Response from Decompose: %v", msg.GetMultiplier())
	}
}
