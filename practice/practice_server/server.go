package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

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

func (*server) LongPractice(stream practicepb.PracticeService_LongPracticeServer) error {
	fmt.Printf("LongPractice functions was invoked with a streaming request")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&practicepb.LongPracticeResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		firstState := req.GetPracticing().GetFirstState()
		result += "The state is " + firstState + "! "
	}
}

func (*server) PracticeBiDi(stream practicepb.PracticeService_PracticeBiDiServer) error {
	fmt.Printf("PracticeBiDi functions was invoked with a streaming request")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}
		firstState := req.GetPracticing().GetFirstState()
		result := "The state is " + firstState + "\n"
		err = stream.Send(&practicepb.PracticeBiDiResponse{
			Result: result,
		})
		if err != nil {
			log.Fatalf("Error while sending data to client stream: %v", err)
			return err
		}
	}
}

func (*server) PracticeWithDeadline(ctx context.Context, req *practicepb.PracticeWithDeadlineRequest) (*practicepb.PracticeWithDeadlineResponse, error) {
	fmt.Printf("PracticeWithDeadline function was invoked with %v\n", req)
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Println("The client canceled the request!")
			return nil, status.Error(codes.Canceled, "The client canceled the request")
		}
		time.Sleep(1 * time.Second)
	}
	firstState := req.GetPracticing().GetFirstState()
	result := "The first state is " + firstState
	res := &practicepb.PracticeWithDeadlineResponse{
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

	tls := false
	opts := []grpc.ServerOption{}
	if tls {
		certFile := ""
		keyFile := ""
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatalf("Failed loading certificates: %v", sslErr)
			return
		}
		opts = append(opts, grpc.Creds(creds))
	}

	s := grpc.NewServer(opts...)
	practicepb.RegisterPracticeServiceServer(s, &server{})

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
