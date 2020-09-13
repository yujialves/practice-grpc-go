package main

import (
	"context"
	"log"

	"../sumpb"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer cc.Close()

	client := sumpb.NewSumServiceClient(cc)
	doUnary(client)
}

func doUnary(client sumpb.SumServiceClient) {

	var firstNumber int64 = 123
	var secondNumber int64 = 456

	req := &sumpb.SumRequest{
		SumData: &sumpb.SumData{
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
