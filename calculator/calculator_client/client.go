package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"../calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer cc.Close()

	client := calculatorpb.NewCalculatorServiceClient(cc)
	// doUnary(client)
	// doServerStreaming(client)
	// doClientStreaming(client)
	// doBiDiStreaming(client)
	doErrorUnary(client)
}

func doUnary(client calculatorpb.CalculatorServiceClient) {

	var firstNumber int64 = 123
	var secondNumber int64 = 456

	req := &calculatorpb.SumRequest{
		SumData: &calculatorpb.SumData{
			FirstNumber:  firstNumber,
			SecondNumber: secondNumber,
		},
	}
	res, err := client.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}
	log.Printf("Response from Sum: %d + %d = %d", firstNumber, secondNumber, res.Result)
}

func doServerStreaming(client calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		PrimeNumber: 120,
	}
	resStream, err := client.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeNumberDecomposition RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from PrimeNumberDecomposition: %d", msg.GetResult())
	}
}

func doClientStreaming(client calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	requests := []*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
			Number: 1,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 2,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 3,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 4,
		},
	}

	stream, err := client.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while calling ComputeAverage: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from ComputeAverage: %v", err)
	}
	fmt.Printf("ComputeAverage Response: %v\n", res)
}

func doBiDiStreaming(client calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a BiDi Streaming RPC...")

	stream, err := client.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	requests := []*calculatorpb.FindMaximumRequest{
		&calculatorpb.FindMaximumRequest{
			Number: -5,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 7,
		},
		&calculatorpb.FindMaximumRequest{
			Number: -10,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 2,
		},
		&calculatorpb.FindMaximumRequest{
			Number: -1,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 9,
		},
	}

	waitc := make(chan struct{})

	go func() {
		for _, req := range requests {
			fmt.Printf("Sending number: %v\n", req)
			stream.Send(req)
		}
		stream.CloseSend()
	}()

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
			fmt.Printf("Received: %v\n", res.GetCurrentMax())
		}
	}()

	<-waitc
}

func doErrorUnary(client calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a SquareRoot Unary RPC...")
	doErrorCall(client, 10)
	doErrorCall(client, -2)
}

func doErrorCall(client calculatorpb.CalculatorServiceClient, n int64) {
	res, err := client.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: n})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			fmt.Printf("Error message from server: %v\n", respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We prabably sent a negative number!")
				return
			}
		} else {
			log.Fatalf("Big error calling SquareRoot: %v", err)
			return
		}
	}
	fmt.Printf("Result of square root of %v: %v\n", n, res.GetNumberRoot())
}
