package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"../practicepb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Practice(ctx context.Context, req *practicepb.PracticeRequest) (*practicepb.PracticeResponse, error) {
	fmt.Printf("Practice function was invoked with %v\n", req)
	firstState := req.GetPracticing().GetFirstState()
	result := "The first state is " + firstState
	res := &practicepb.PracticeResponse{
		Result: result,
	}
	return res, nil
}

func (*server) PracticeManyTimes(req *practicepb.PracticeManyTimesRequest, stream practicepb.PracticeService_PracticeManyTimesServer) error {
	fmt.Printf("PracticeManyTimes functions was invoked with %v\n", req)
	firstState := req.GetPracticing().GetFirstState()
	for i := 0; i < 10; i++ {
		result := "FirstState: " + firstState + ", Number: " + strconv.Itoa(i)
		res := &practicepb.PracticeManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
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
