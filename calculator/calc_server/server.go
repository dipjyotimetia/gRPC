package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/dipjyotimetia/gogrpc/calculator/calcpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calcpb.SumRequest) (*calcpb.SumResponse, error) {
	fmt.Printf("Sum function was invoked %v\n", req)
	firstArgument := req.GetFirstNumber()
	secondArgument := req.GetSecondNumber()

	result := firstArgument + secondArgument
	res := &calcpb.SumResponse{
		Result: result,
	}
	return res, nil
}

func main() {

	fmt.Println("Hello calc")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen server")
	}

	s := grpc.NewServer()
	calcpb.RegisterSumServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
