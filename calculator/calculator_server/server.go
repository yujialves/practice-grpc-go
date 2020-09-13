package main

import (
	"context"
	"fmt"
	"io"
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

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverage functions was invoked with a streaming request")
	var sum int64
	var cnt int64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Result: float64(sum) / float64(cnt),
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		num := req.GetNumber()
		sum += num
		cnt++
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("FindMaximum functions was invoked with a streaming request")

	var max int64
	for i := 0; ; i++ {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}
		if i == 0 {
			max = req.GetNumber()
		} else {
			if req.GetNumber() > max {
				max = req.GetNumber()
			}
		}

		err = stream.Send(&calculatorpb.FindMaximumResponse{
			CurrentMax: max,
		})
		if err != nil {
			log.Fatalf("Error while sending data to client stream: %v", err)
			return err
		}
	}
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
