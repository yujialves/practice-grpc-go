package main

import (
	"context"
	"log"
	"net"

	"../calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	sumData := req.GetSumData()
	result := sumData.GetFirstNumber() + sumData.GetSecondNumber()
	res := &calculatorpb.SumResponse{
		Result: result,
	}
	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	primeNumber := req.GetPrimeNumber()
	var num int64 = 2
	for primeNumber > 1 {
		if primeNumber%num == 0 {
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				Result: num,
			}
			stream.Send(res)
			primeNumber /= num
		} else {
			num++
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
