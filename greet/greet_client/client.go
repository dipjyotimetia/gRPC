package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/dipjyotimetia/gogrpc/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello i am a client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("error is, %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	// doUnary(c)
	// doServerStreaming(c)
	doClientStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Staring to do unary rpc")
	req := &greetpb.GreetRequest{Greeting: &greetpb.Greeting{
		FirstName: "Dipjyoti",
		LastName:  "Metia",
	}}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling greet rpc %v", err)
	}
	log.Printf("Response from greet %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Dipjyoti",
			LastName:  "Metia",
		}}
	res, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling greet rpc %v", err)
	}
	for {
		msg, err := res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading the stream %v", err)
		}
		log.Printf("Response from greet many times %v", msg.GetResult())
	}

}

func doClientStreaming(c greetpb.GreetServiceClient) {

	request := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Dipjyoti1",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Dipjyoti2",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Dipjyoti3",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Dipjyoti4",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Dipjyoti5",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Dipjyoti6",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())

	if err != nil {
		log.Fatalf("Error while reading")
	}

	for _, req := range request {
		fmt.Printf("Sending request %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}
	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from long greet: %v", err)
	}
	fmt.Printf("LongGreet response being %v", response)
}
