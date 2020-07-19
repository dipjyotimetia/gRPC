package main

import (
	"context"
	"fmt"
	"log"

	"github.com/dipjyotimetia/gogrpc/calculator/calcpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {

	fmt.Println("Hello i am a calc client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("error is %v", err)
	}

	defer cc.Close()

	c := calcpb.NewSumServiceClient(cc)
	// doUnary(c)
	doErrorUnary(c, 10)
	doErrorUnary(c, -2)
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

func doErrorUnary(c calcpb.SumServiceClient, number int32) {
	fmt.Println("Starting to do a sqrt call")

	res, err := c.SquareRoot(context.Background(), &calcpb.SquareRootRequest{Number: number})

	if err != nil {
		resErr, ok := status.FromError(err)
		if ok {
			fmt.Println(resErr.Message())
			fmt.Println(resErr.Code())
			if resErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative error")
			}
		} else {
			log.Fatalf("Big error calling squareroot: %v", err)
		}

	}
	fmt.Printf("Result of square root of %v", res)
}
