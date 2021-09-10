package main

import (
	"context"
	"fmt"
	"grpc-g-course/greet/greet/greetpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main(){
	fmt.Println("Hi from client")
	cc,err:=grpc.Dial("localhost:50051",grpc.WithInsecure())
	if err!=nil{
		log.Fatalf("failed to connect: %v",err)
	}
    defer cc.Close()  //to defer the execution to end of program
	
	c:=greetpb.NewGreetServiceClient(cc)
	//fmt.Println("created client : %f",c)
	//doUnary(c)
	doServerStreaming(c)

}
func doUnary(c greetpb.GreetServiceClient){
	fmt.Println("starting to do unary RPC")
	req:=&greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{FirstName: "Sreeroop", LastName: "Shiv"},
	}

	res,err:= c.Greet(context.Background(), req)
		if err!=nil{
			log.Fatalf("Error while calling Greet RPC : %v",err)
		}
		log.Printf("Response from Greet: %v",res.Result )

	}

func doServerStreaming(c greetpb.GreetServiceClient){
	fmt.Println("starting to do server streaming RPC")
	req:=&greetpb.GreeetManyRequest{
		Greeting: &greetpb.Greeting{FirstName: "Sreeroop Shiv",
	LastName: "UK",},
	}

	resStream,err:=c.GreetManyTimes(context.Background(),req)
	if err!=nil{
		log.Fatalf("Error while calling GreetMany RPC: %v",err)
	}
	for { //to accept as many message 
		msg,err:=resStream.Recv()

	if err==io.EOF{
		//we have reached end of stream
		break
	}
	if err!=nil{
		log.Fatalf("error while reading stream : %v",err)
	}
	log.Printf("Streamed message is : %v",msg.Result)
}
}	


