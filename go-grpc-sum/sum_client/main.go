package main

import (
	"context"
	"log"
	"time"

	sum "github.com/ReigoHein/udemy-grpc-course/go-grpc-sum/models"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := sum.NewSumClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Add(ctx, &sum.AddRequest{Numbers: []int32{1, 2, 3, 4, 5}})
	if err != nil {
		log.Fatalf("could not add: %v", err)
	}
	log.Printf("Result: %d", r.Result)
}
