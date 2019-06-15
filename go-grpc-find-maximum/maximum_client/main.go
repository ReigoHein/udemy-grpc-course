package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	maximum "github.com/ReigoHein/udemy-grpc-course/go-grpc-find-maximum/models"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Find maximum client started")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := maximum.NewFindMaximumServiceClient(cc)

	findMaximum(c)
}

func findMaximum(c maximum.FindMaximumServiceClient) {
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	requests := []*maximum.MaximumRequest{}
	testData := []int{1, 5, 3, 6, 2, 20, 2, 20, 20, 50}
	for _, val := range testData {
		maximumRequest := &maximum.MaximumRequest{
			Number: &maximum.Number{
				Number: int32(val),
			},
		}
		requests = append(requests, maximumRequest)
	}

	fmt.Printf("Requests: %v, len %v\n", requests, len(requests))

	waitc := make(chan struct{})

	// send messages to the server
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending number: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// receive messages from the server
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("Maximum: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	<-waitc
}
