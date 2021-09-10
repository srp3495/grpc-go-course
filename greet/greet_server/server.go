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

func(*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error{
	res:="Hello "
	for{
	req,err:=stream.Recv()
	

	if err==io.EOF{
		fmt.Printf("End of file reached")
		return stream.SendAndClose(&greetpb.LongGreetResponse{
			Result: res,
		})
	}
	if err!=nil{
		log.Fatalf("Error receiving client stream : %v",err)
	}
    first:=req.Greeting.FirstName
	res+=first+" !"
	
}
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
