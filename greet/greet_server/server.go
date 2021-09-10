package main

import (
	"context"
	"fmt"
	"grpc-g-course/greet/greet/greetpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer //why i am not sure
}

func (*server) Greet(ctx context.Context, in *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function is invoke fro req: %v",in)
	first_name := in.Greeting.FirstName
	result := "Hello " + first_name
	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil

}

func main() {

	fmt.Println("Hi Server")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("Failed to listen", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
