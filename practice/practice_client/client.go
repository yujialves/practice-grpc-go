package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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
	// fmt.Printf("Created client: %f", client)
	// doUnary(client)
	// doServerStreaming(client)
	// doClientStreaming(client)
	doBiDiStreaming(client)
}

func doUnary(client practicepb.PracticeServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &practicepb.PracticeRequest{
		Practicing: &practicepb.Practicing{
			FirstState:  "Hello",
			SecondState: "World",
		},
	}
	res, err := client.Practice(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Practice RPC: %v", err)
	}
	log.Printf("Response from Practice: %v", res.Result)
}

func doServerStreaming(client practicepb.PracticeServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")

	req := &practicepb.PracticeManyTimesRequest{
		Practicing: &practicepb.Practicing{
			FirstState:  "Hello",
			SecondState: "World",
		},
	}
	resStream, err := client.PracticeManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PracticeManyTimes RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from PracticeManyTimes: %v", msg.GetResult())
	}
}

func doClientStreaming(client practicepb.PracticeServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	requests := []*practicepb.LongPracticeRequest{
		&practicepb.LongPracticeRequest{
			Practicing: &practicepb.Practicing{
				FirstState: "Hello",
			},
		},
		&practicepb.LongPracticeRequest{
			Practicing: &practicepb.Practicing{
				FirstState: "World",
			},
		},
		&practicepb.LongPracticeRequest{
			Practicing: &practicepb.Practicing{
				FirstState: "hello",
			},
		},
		&practicepb.LongPracticeRequest{
			Practicing: &practicepb.Practicing{
				FirstState: "world",
			},
		},
		&practicepb.LongPracticeRequest{
			Practicing: &practicepb.Practicing{
				FirstState: "Hello World",
			},
		},
		&practicepb.LongPracticeRequest{
			Practicing: &practicepb.Practicing{
				FirstState: "hello world",
			},
		},
	}

	stream, err := client.LongPractice(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongPractice: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from LongPractice: %v", err)
	}
	fmt.Printf("LongPractice Response: %v\n", res)
}

func doBiDiStreaming(client practicepb.PracticeServiceClient) {
	fmt.Println("Starting to do a BiDi Streaming RPC...")

	stream, err := client.PracticeBiDi(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	requests := []*practicepb.PracticeBiDiRequest{
		&practicepb.PracticeBiDiRequest{
			Practicing: &practicepb.Practicing{
				FirstState: "Hello",
			},
		},
		&practicepb.PracticeBiDiRequest{
			Practicing: &practicepb.Practicing{
				FirstState: "World",
			},
		},
		&practicepb.PracticeBiDiRequest{
			Practicing: &practicepb.Practicing{
				FirstState: "hello",
			},
		},
		&practicepb.PracticeBiDiRequest{
			Practicing: &practicepb.Practicing{
				FirstState: "world",
			},
		},
		&practicepb.PracticeBiDiRequest{
			Practicing: &practicepb.Practicing{
				FirstState: "Hello World",
			},
		},
		&practicepb.PracticeBiDiRequest{
			Practicing: &practicepb.Practicing{
				FirstState: "hello world",
			},
		},
	}

	waitc := make(chan struct{})

	go func() {
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
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
			fmt.Printf("Received: %v\n", res.GetResult())
		}
	}()

	<-waitc
}
