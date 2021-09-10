package main

import (
	"context"
	"fmt"
	"grpc-g-course/greet/greet/greetpb"
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
	doUnary(c)

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


