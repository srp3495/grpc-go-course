package main

import (
	"context"
	"fmt"
	"grpc-g-course/greet/calculator/calculatorpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Calculator client")
	cc, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer cc.Close() //to defer the execution to end of program

	c := calculatorpb.NewCalculatorServiceClient(cc)
	//fmt.Println("created client : %f",c)
	//doUnary(c)
	doStreaming(c)

}
func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("starting to do unary RPC")
	req := &calculatorpb.SumRequest{
		FirstNum:  5,
		SecondNum: 8,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC : %v", err)
	}
	log.Printf("Response from Sum: %v", res.Result)

}

func doStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("starting to do streaming server RPC")
	req := &calculatorpb.PrimeDecompositionStreamRequest{
		Value: 100,
	}
	res, err := c.PrimeDecompositionStream(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling streaming rpc : %v", err)
	}
	for {
		msg, err := res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened : %v", err)
		}
		fmt.Printf("Received number is: %v \n", msg.Result)
		time.Sleep(1000 * time.Millisecond)
	}
}
