package main

import (
	"context"
	"fmt"
	"log"
	"time"

	average "github.com/ReigoHein/udemy-grpc-course/go-grpc-compute-average/models"
	"google.golang.org/grpc"
)

type server struct{}

func main() {
	fmt.Println("Started compute average client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := average.NewComputeAverageServiceClient(cc)

	calculateAverages(c)
}

func calculateAverages(c average.ComputeAverageServiceClient) {
	fmt.Println("Client streaming to compute average RPC")

	requests := []*average.AverageRequest{
		&average.AverageRequest{
			Average: &average.Average{
				Number: 5,
			},
		},
		&average.AverageRequest{
			Average: &average.Average{
				Number: 22,
			},
		},
		&average.AverageRequest{
			Average: &average.Average{
				Number: 37,
			},
		},
		&average.AverageRequest{
			Average: &average.Average{
				Number: 32,
			},
		},
		&average.AverageRequest{
			Average: &average.Average{
				Number: 10,
			},
		},
	}

	stream, err := c.ComputeAverage(context.Background())

	if err != nil {
		log.Fatalf("Error while calling ComputeAverage: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending request: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from ComputeAverage: %v", err)
	}

	fmt.Printf("ComputeAverage response: %+v\n", res)
}
