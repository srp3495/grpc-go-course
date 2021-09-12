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
	//doStreaming(c)
	//doClientStreaming(c)
    doBidiStreaming(c)

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

func doClientStreaming(c calculatorpb.CalculatorServiceClient){
	request:=[]*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
			Value: 30,
		},
		&calculatorpb.ComputeAverageRequest{
			Value: 40,
		},
		&calculatorpb.ComputeAverageRequest{
			Value: 50,
		},
		&calculatorpb.ComputeAverageRequest{
			Value: 60,
		},
	}

    
	stream,err:=c.ComputeAverage(context.Background())
	if err!=nil{
		log.Fatalf("Error streaming :%v",err)

	}

	for _,req:=range request{
        fmt.Printf("\n sending the request: %v",req.Value)
		stream.Send(req)
		time.Sleep(1000*time.Millisecond)
	}
	res,err:=stream.CloseAndRecv()
	if err!=nil{
		log.Fatalf("Error in receiving the response")
	}
	fmt.Printf("\n Received result: %v",res)

}

func doBidiStreaming(c calculatorpb.CalculatorServiceClient){
	request:=[]int32{2,5,7,4,9,12,4,67,66,87,12}

	waitc:=make(chan struct{})

	stream,err:=c.SendMaximum(context.Background())

	if err!=nil{
		log.Fatalf("Error connecting to stream : %v",err)
		return
	}

	//send bunch of messages
	go func ()  {
		for _,req:=range request{
			fmt.Printf("\n Sending number : %v",req)
			stream.Send(&calculatorpb.MaxRequest{
				Value: req,
			})
			time.Sleep(1000*time.Millisecond)
		}
		stream.CloseSend()
	}()

    //receive bunch of messages
	go func ()  {
		for{
			res,err:=stream.Recv()
			if err==io.EOF {
				fmt.Print("\n End of message reached \n")
				break
			}
			if err!=nil{
				log.Fatalf("Error receiving stream : %v",err)
				break
			}

             fmt.Printf("\n Received response : %v",res.Result)
			}
			close(waitc)
	}()
	
	<-waitc //to come here when channel is closed

}
