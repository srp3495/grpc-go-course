package main

import (
	"context"
	"fmt"
	"grpc-g-course/greet/greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hi from client")
	tls := true //to disable ssl handshake, for enabling make it true
	opts := grpc.WithInsecure()
	if tls {
		certFile := "ssl/ca.crt" //CA trust certificate
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Error loading certificate : %v", sslErr)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	//cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer cc.Close() //to defer the execution to end of program

	c := greetpb.NewGreetServiceClient(cc)
	//fmt.Println("created client : %f",c)
	doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	//dpBiDiStreaming(c)
	//doUnaryWithDeadline(c, 5*time.Second) //should complete
	//doUnaryWithDeadline(c, 1*time.Second) //should timeout

}
func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("starting to do unary RPC")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{FirstName: "Sreeroop", LastName: "Shiv"},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC : %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)

}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("starting to do server streaming RPC")
	req := &greetpb.GreeetManyRequest{
		Greeting: &greetpb.Greeting{FirstName: "Sreeroop Shiv",
			LastName: "UK"},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetMany RPC: %v", err)
	}
	for { //to accept as many message
		msg, err := resStream.Recv()

		if err == io.EOF {
			//we have reached end of stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream : %v", err)
		}
		log.Printf("Streamed message is : %v", msg.Result)
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("starting to do client streaming RPC")

	request := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Sreeroop",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Sreeshma",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Sivadasan",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Sudarsini",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error streaming client : %v", err)
	}

	for _, req := range request {
		fmt.Printf("\nSending request : %v", req.Greeting.FirstName)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving the result from server : %v", err)
	}

	fmt.Printf("\n Recieved result is : %v", res)

}

func dpBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("starting Bidi streaming")
	//make request
	request := []*greetpb.GreetEveryOneRequest{
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Sreeroop",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Sreeshma",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Sivadasan",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{FirstName: "Sudarsini"},
		},
	}

	//create a stream by invoking client
	stream, err := c.GreetEveryOne(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream :%v", err)
		return
	}

	//for waiting
	waitc := make(chan struct{})

	//send bunch of message to the server(go routine)
	go func() {
		for _, req := range request {
			fmt.Printf("\n Sending message: %v", req.Greeting.FirstName)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend() //after sending all messages
	}()

	//reveive a bunch of messages from server

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				fmt.Print("\n End of message reached \n")
				break
			}
			if err != nil {
				log.Fatalf("Error receiving stream : %v", err)
				break
			}

			fmt.Printf("\n Received response : %v", res.Result)
		}
		close(waitc)
	}()

	<-waitc //to come here when channel is closed
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("starting to do unary with deadline RPC")
	req := &greetpb.DeadLineGreetRequest{
		Greeting: &greetpb.Greeting{FirstName: "Sreeroop", LastName: "Shiv"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)

	if err != nil {

		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Deadline exceeded ")
			} else {
				fmt.Printf("some framework error happened : %v \n", err)
			}
		} else {
			log.Fatalf("Error while calling Greet RPC : %v \n", err)
		}
		return
	}
	log.Printf("Response from Greet: %v \n", res.Result)

}
