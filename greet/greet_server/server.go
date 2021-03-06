package main

import (
	"context"
	"fmt"
	"grpc-g-course/greet/greet/greetpb"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer //why i am not sure
}

func (*server) Greet(ctx context.Context, in *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function is invoke fro req: %v", in)
	first_name := in.Greeting.FirstName
	result := "Hello " + first_name
	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil

}

func (*server) GreetManyTimes(req *greetpb.GreeetManyRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function is invoke fro req: %v", req)
	first := req.Greeting.FirstName

	for i := 0; i < 10; i++ {
		result := "Hello " + first + " numbered :" + strconv.Itoa(i)
		res := &greetpb.GreeetManyResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)

	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet function is invoked")
	res := "Hello "
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			fmt.Printf("End of file reached")
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: res,
			})
		}
		if err != nil {
			log.Fatalf("Error receiving client stream : %v", err)
		}
		first := req.Greeting.FirstName
		res += first + " !"

	}
}

func (*server) GreetEveryOne(stream greetpb.GreetService_GreetEveryOneServer) error {
	fmt.Printf("GreetEveryOne function is invoked")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error receiving stream : %v", err)
		}

		first := req.Greeting.FirstName
		result := "Hello " + first + " !"

		Senderr := stream.Send(&greetpb.GreetEveryOneResponse{
			Result: result,
		})
		if Senderr != nil {
			log.Fatalf("Error sending response to client : %v", Senderr)
			return Senderr
		}

	}
}

func (*server) GreetWithDeadline(ctx context.Context, in *greetpb.DeadLineGreetRequest) (*greetpb.DeadLineGreetResponse, error) {
	fmt.Printf("Greet function with deadline is invoke from req: %v", in)

	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			//client cancelled request
			fmt.Println("Deadline exceeded, client cancelled request")
			return nil, status.Error(codes.DeadlineExceeded, "client cancelled request")
		}
		time.Sleep(1 * time.Second)
	}

	first_name := in.Greeting.GetFirstName()
	result := "Hello " + first_name
	res := &greetpb.DeadLineGreetResponse{
		Result: result,
	}

	return res, nil

}

func main() {

	fmt.Println("Hi Server")

	tls := false //to disable ssl handshake, for enabling make it true
	opts := []grpc.ServerOption{}
	if tls {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatalf("Error loading ssl certificate : %v", sslErr)
		}

		opts = append(opts, grpc.Creds(creds))
	} else {
		opts = nil
	}

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("Failed to listen", err)
	}

	s := grpc.NewServer(opts...)
	greetpb.RegisterGreetServiceServer(s, &server{})

	reflection.Register(s)


	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
