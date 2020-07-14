package main

import (
	"context"
	"fmt"
	"log"

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

	doUnary(c)
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
