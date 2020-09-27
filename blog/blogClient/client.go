package main

import (
	"context"
	"fmt"
	blogpb "github.com/dipjyotimetia/gogrpc/blog/blogPb"
	"google.golang.org/grpc"
	"log"
)

func main() {
	fmt.Println("Hello am a blogging client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)

	if err != nil {
		log.Fatalf("error is, %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	req := &blogpb.Blog{
		AuthorId: "Dip",
		Title:    "Dip gRPC Mongo",
		Content:  "mongo service",
	}
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: req})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Blog has been created", res)
}
