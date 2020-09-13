package main

import (
	"context"
	"log"

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
	doUnary(client)
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
