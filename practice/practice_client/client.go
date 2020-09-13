package main

import (
	"fmt"
	"log"

	"../practicepb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello World from client.")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer cc.Close()

	client := practicepb.NewPracticeServiceClient(cc)
	fmt.Printf("Created client: %f", client)
}
