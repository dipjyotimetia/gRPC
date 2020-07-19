package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/dipjyotimetia/gogrpc/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	// doClientStreaming(c)
	// doClientBidirectional(c)

	doUnaryWithDeadline(c, 5*time.Second)
	doUnaryWithDeadline(c, 1*time.Second)
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
	fmt.Println("Starting to do server streaming rpc")

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
	fmt.Println("Starting to do client streaming rpc")

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

func doClientBidirectional(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do bidirectional rpc")

	request := []*greetpb.GreetEveryoneRequest{
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

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while streaming")
		return
	}

	waitC := make(chan struct{})

	go func() {
		for _, req := range request {
			fmt.Printf("Sending message: %v", req)
			stream.Send(req)
			time.Sleep(100 * time.Millisecond)
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
				log.Fatalf("Error while receiving %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitC)
	}()

	<-waitC
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, times time.Duration) {
	fmt.Println("Staring to do unary with deadline rpc")

	req := &greetpb.GreetWithDeadLineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Dipjyoti",
			LastName:  "Metia",
		}}

	ctx, cancel := context.WithTimeout(context.Background(), times)

	defer cancel()

	res, err := c.GreetWithDeadLine(ctx, req)
	if err != nil {
		statusError, ok := status.FromError(err)
		if ok {
			if statusError.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline was exceeded")
			} else {
				fmt.Printf("Unexpected error: %v", err)
			}
		} else {
			log.Fatalf("Error while calling GreetWIthDeadline RPC:%v", err)
		}
		return
	}
	log.Printf("Response from greet with deadline: %v", res.Result)

}
