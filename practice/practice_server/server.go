package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"../practicepb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Practice(ctx context.Context, req *practicepb.PracticeRequest) (*practicepb.PracticeResponse, error) {
	firstState := req.GetPracticing().GetFirstState()
	result := "The first state is" + firstState
	res := &practicepb.PracticeResponse{
		Result: result,
	}
	return res, nil
}

func main() {
	fmt.Println("Hello World")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	practicepb.RegisterPracticeServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
