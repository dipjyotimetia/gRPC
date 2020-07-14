package main

import (
	"context"
	"fmt"
	"log"

	"github.com/dipjyotimetia/gogrpc/calculator/calcpb"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Hello i am a calc client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("error is %v", err)
	}

	defer cc.Close()

	c := calcpb.NewSumServiceClient(cc)
	doUnary(c)
}

func doUnary(c calcpb.SumServiceClient) {
	fmt.Println("Starting to do unary rpc")

	req := &calcpb.SumRequest{FirstNumber: 2, SecondNumber: 4}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling calc rpc %v", err)
	}
	log.Printf("response from sum %v", res)
}
